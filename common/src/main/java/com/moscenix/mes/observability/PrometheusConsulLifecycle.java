package com.moscenix.mes.observability;

import com.moscenix.mes.consul.ConsulServiceRegistrar;
import java.util.List;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.SmartLifecycle;
import org.springframework.stereotype.Component;

@Component
public class PrometheusConsulLifecycle implements SmartLifecycle {
    private final boolean enabled;
    private final ConsulServiceRegistrar registrar;
    private boolean running;

    public PrometheusConsulLifecycle(
            @Value("${observability.prometheus.consul-registration-enabled:true}") boolean enabled,
            @Value("${consul.address:}") String consulAddress,
            @Value("${observability.prometheus.service-name:prometheus}") String prometheusServiceName,
            @Value("${observability.prometheus.advertise-host:${grpc.server.advertise-host:127.0.0.1}}") String advertiseHost,
            @Value("${grpc.server.service-name:${spring.application.name}}") String serviceName,
            @Value("${server.port}") int metricsPort) {
        this.enabled = enabled;
        this.registrar = new ConsulServiceRegistrar(
                consulAddress,
                prometheusServiceName + "-" + serviceName + "-" + metricsPort,
                prometheusServiceName,
                advertiseHost,
                metricsPort,
                List.of("service:" + serviceName));
    }

    @Override
    public void start() {
        if (enabled) {
            registrar.register();
        }
        running = true;
    }

    @Override
    public void stop() {
        if (enabled) {
            registrar.deregister();
        }
        running = false;
    }

    @Override
    public boolean isRunning() {
        return running;
    }
}
