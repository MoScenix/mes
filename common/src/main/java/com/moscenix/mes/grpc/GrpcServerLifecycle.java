package com.moscenix.mes.grpc;

import com.moscenix.mes.consul.ConsulServiceRegistrar;
import io.grpc.BindableService;
import io.grpc.Server;
import io.grpc.ServerBuilder;
import io.grpc.ServerInterceptor;
import io.grpc.ServerInterceptors;
import java.io.IOException;
import java.util.List;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.SmartLifecycle;
import org.springframework.stereotype.Component;

@Component
public class GrpcServerLifecycle implements SmartLifecycle {
    private final Server server;
    private final ConsulServiceRegistrar consulRegistrar;
    private Thread awaitThread;
    private boolean running;

    public GrpcServerLifecycle(
            List<BindableService> services,
            List<ServerInterceptor> interceptors,
            @Value("${grpc.server.port}") int port,
            @Value("${grpc.server.service-name}") String serviceName,
            @Value("${grpc.server.advertise-host:127.0.0.1}") String advertiseHost,
            @Value("${consul.address:}") String consulAddress) {
        ServerBuilder<?> builder = ServerBuilder.forPort(port);
        services.forEach(service -> builder.addService(interceptors.isEmpty()
                ? service.bindService()
                : ServerInterceptors.intercept(service.bindService(), interceptors)));
        this.server = builder.build();
        this.consulRegistrar = new ConsulServiceRegistrar(
                consulAddress, serviceName + "-" + port, serviceName, advertiseHost, port, List.of());
    }

    @Override
    public void start() {
        try {
            server.start();
            consulRegistrar.register();
            awaitThread = new Thread(this::awaitTermination, "grpc-server-await");
            awaitThread.setDaemon(false);
            awaitThread.start();
            running = true;
        } catch (IOException e) {
            throw new IllegalStateException("failed to start grpc server", e);
        }
    }

    @Override
    public void stop() {
        consulRegistrar.deregister();
        server.shutdown();
        if (awaitThread != null) {
            awaitThread.interrupt();
            awaitThread = null;
        }
        running = false;
    }

    @Override
    public boolean isRunning() {
        return running;
    }

    private void awaitTermination() {
        try {
            server.awaitTermination();
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
    }

}
