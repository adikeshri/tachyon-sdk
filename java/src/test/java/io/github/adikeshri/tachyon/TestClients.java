package io.github.adikeshri.tachyon;

final class TestClients {
    private TestClients() {
    }

    static TachyonClient forServer(TestServer server) {
        return TachyonClient.builder(server.baseUrl).apiKey("admin-key").build();
    }
}
