package com.moscenix.mes.consul;

import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.util.List;
import java.util.stream.Collectors;

public class ConsulServiceRegistrar {
    private final HttpClient httpClient = HttpClient.newHttpClient();
    private final String consulAddress;
    private final String serviceId;
    private final String serviceName;
    private final String advertiseHost;
    private final int port;
    private final List<String> tags;

    public ConsulServiceRegistrar(
            String consulAddress,
            String serviceId,
            String serviceName,
            String advertiseHost,
            int port,
            List<String> tags) {
        this.consulAddress = normalize(consulAddress);
        this.serviceId = serviceId;
        this.serviceName = serviceName;
        this.advertiseHost = advertiseHost;
        this.port = port;
        this.tags = tags == null ? List.of() : List.copyOf(tags);
    }

    public void register() {
        if (consulAddress.isEmpty()) {
            return;
        }
        String tagsJson = tags.stream()
                .map(tag -> "\"" + escape(tag) + "\"")
                .collect(Collectors.joining(",", "[", "]"));
        String body = """
                {"ID":"%s","Name":"%s","Address":"%s","Port":%d,"Tags":%s}
                """.formatted(escape(serviceId), escape(serviceName), escape(advertiseHost), port, tagsJson);
        send("PUT", "/v1/agent/service/register", body);
    }

    public void deregister() {
        if (consulAddress.isEmpty()) {
            return;
        }
        send("PUT", "/v1/agent/service/deregister/" + serviceId, "");
    }

    private void send(String method, String path, String body) {
        try {
            HttpRequest request = HttpRequest.newBuilder(URI.create(consulAddress + path))
                    .method(method, HttpRequest.BodyPublishers.ofString(body))
                    .header("Content-Type", "application/json")
                    .build();
            httpClient.send(request, HttpResponse.BodyHandlers.discarding());
        } catch (IOException e) {
            throw new IllegalStateException("failed to call consul", e);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new IllegalStateException("interrupted while calling consul", e);
        }
    }

    private static String normalize(String address) {
        if (address == null || address.isBlank()) {
            return "";
        }
        String trimmed = address.trim();
        if (trimmed.startsWith("http://") || trimmed.startsWith("https://")) {
            return trimmed;
        }
        return "http://" + trimmed;
    }

    private static String escape(String value) {
        if (value == null) {
            return "";
        }
        return value.replace("\\", "\\\\").replace("\"", "\\\"");
    }
}
