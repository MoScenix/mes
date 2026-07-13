package com.moscenix.mes.observability;

import io.grpc.ForwardingServerCall;
import io.grpc.Metadata;
import io.grpc.ServerCall;
import io.grpc.ServerCallHandler;
import io.grpc.ServerInterceptor;
import io.grpc.Status;
import io.micrometer.core.instrument.MeterRegistry;
import io.micrometer.core.instrument.Timer;
import java.util.concurrent.atomic.AtomicBoolean;
import org.springframework.stereotype.Component;

@Component
public class GrpcMetricsServerInterceptor implements ServerInterceptor {
    private final MeterRegistry meterRegistry;

    public GrpcMetricsServerInterceptor(MeterRegistry meterRegistry) {
        this.meterRegistry = meterRegistry;
    }

    @Override
    public <ReqT, RespT> ServerCall.Listener<ReqT> interceptCall(
            ServerCall<ReqT, RespT> call, Metadata headers, ServerCallHandler<ReqT, RespT> next) {
        String fullMethodName = call.getMethodDescriptor().getFullMethodName();
        String grpcService = grpcService(fullMethodName);
        String grpcMethod = grpcMethod(fullMethodName);
        Timer.Sample sample = Timer.start(meterRegistry);
        AtomicBoolean recorded = new AtomicBoolean(false);

        ServerCall<ReqT, RespT> monitoringCall = new ForwardingServerCall.SimpleForwardingServerCall<>(call) {
            @Override
            public void close(Status status, Metadata trailers) {
                recordOnce(sample, recorded, grpcService, grpcMethod, status);
                super.close(status, trailers);
            }
        };

        try {
            return next.startCall(monitoringCall, headers);
        } catch (RuntimeException e) {
            recordOnce(sample, recorded, grpcService, grpcMethod, Status.fromThrowable(e));
            throw e;
        }
    }

    private void recordOnce(
            Timer.Sample sample,
            AtomicBoolean recorded,
            String grpcService,
            String grpcMethod,
            Status status) {
        if (!recorded.compareAndSet(false, true)) {
            return;
        }
        sample.stop(Timer.builder("rpc.server.duration")
                .description("gRPC server call duration")
                .tag("rpc.system", "grpc")
                .tag("rpc.service", grpcService)
                .tag("rpc.method", grpcMethod)
                .tag("rpc.status", status.getCode().name())
                .register(meterRegistry));
    }

    private static String grpcService(String fullMethodName) {
        int slashIndex = fullMethodName.lastIndexOf('/');
        return slashIndex >= 0 ? fullMethodName.substring(0, slashIndex) : fullMethodName;
    }

    private static String grpcMethod(String fullMethodName) {
        int slashIndex = fullMethodName.lastIndexOf('/');
        return slashIndex >= 0 ? fullMethodName.substring(slashIndex + 1) : fullMethodName;
    }
}
