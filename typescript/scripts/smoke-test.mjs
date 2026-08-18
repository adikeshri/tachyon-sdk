// Exercises the *built* dist/ output with plain `node`, no vitest involved.
//
// vitest 4 pulls in `rolldown`, which needs `node:util`'s `styleText` export
// (added in Node 20.12) just to boot — so it can't run at all on Node 18,
// even though the package's own runtime code (this script's target) only
// needs fetch/URL/AbortController, all present since Node 18. This script
// is what actually backs the `engines.node >= 18` claim in package.json: it
// runs on a real Node 18 in CI, importing the compiled output the same way
// a consumer would, with no vitest/rolldown in the loop to get in the way.

import { Tachyon, TachyonError } from '../dist/index.js';

function makeFetch(handler) {
  return async (url) => {
    const { status, body } = handler(String(url));
    return new Response(JSON.stringify(body), {
      status,
      headers: { 'content-type': 'application/json' },
    });
  };
}

async function main() {
  const fetchImpl = makeFetch((url) => {
    if (url.endsWith('/health')) {
      return { status: 200, body: { ok: true, version: 'smoke', uptime_seconds: 1, num_collections: 0 } };
    }
    return { status: 404, body: { error: { code: 'collection_not_found', message: 'no such collection' } } };
  });

  const client = new Tachyon({ url: 'http://smoke-test.invalid', fetch: fetchImpl });

  const health = await client.health();
  if (health.ok !== true || health.version !== 'smoke') {
    throw new Error(`unexpected health response: ${JSON.stringify(health)}`);
  }

  let caught;
  try {
    await client.collections.retrieve('missing');
  } catch (err) {
    caught = err;
  }
  if (!(caught instanceof TachyonError) || caught.code !== 'collection_not_found' || caught.status !== 404) {
    throw new Error(`expected a 404 TachyonError, got: ${caught}`);
  }

  console.log(`smoke test passed on Node ${process.version}`);
}

main().catch((err) => {
  console.error(err);
  process.exitCode = 1;
});
