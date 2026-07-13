package com.moscenix.mes.observability;

import io.grpc.ForwardingServerCall;
import io.grpc.ForwardingServerCallListener;
import io.grpc.Metadata;
import io.grpc.ServerCall;
import io.grpc.ServerCallHandler;
import io.grpc.ServerInterceptor;
import io.grpc.Status;
import io.micrometer.tracing.Span;
import io.micrometer.tracing.Tracer;
import java.util.concurrent.atomic.AtomicBoolean;
import org.springframework.beans.factory.ObjectProvider;
import org.springframework.stereotype.Component;

@Component
public class GrpcTracingServerInterceptor implements ServerInterceptor {
    private final ObjectProvider<Tracer> tracerProvider;

    public GrpcTracingServerInterceptor(ObjectProvider<Tracer> tracerProvider) {
        this.tracerProvider = tracerProvider;
    }

    @Override
    public <ReqT, RespT> ServerCall.Listener<ReqT> interceptCall(
            ServerCall<ReqT, RespT> call, Metadata headers, ServerCallHandler<ReqT, RespT> next) {
        Tracer tracer = tracerProvider.getIfAvailable();
        if (tracer == null) {
            return next.startCall(call, headers);
        }

        String fullMethodName = call.getMethodDescriptor().getFullMethodName();
        String grpcService = grpcService(fullMethodName);
        String grpcMethod = grpcMethod(fullMethodName);
        Span span = tracer.nextSpan()
                .name(fullMethodName)
                .tag("rpc.system", "grpc")
                .tag("rpc.service", grpcService)
                .tag("rpc.method", grpcMethod)
                .start();
        AtomicBoolean ended = new AtomicBoolean(false);

        ServerCall<ReqT, RespT> tracingCall = new ForwardingServerCall.SimpleForwardingServerCall<>(call) {
            @Override
            public void close(Status status, Metadata trailers) {
                span.tag("rpc.status", status.getCode().name());
                if (!status.isOk() && status.getCause() != null) {
                    span.error(status.getCause());
                }
                try {
                    super.close(status, trailers);
                } finally {
                    endOnce(span, ended);
                }
            }
        };

        try (Tracer.SpanInScope scope = tracer.withSpan(span)) {
            ServerCall.Listener<ReqT> listener = next.startCall(tracingCall, headers);
            return new ForwardingServerCallListener.SimpleForwardingServerCallListener<>(listener) {
                @Override
                public void onMessage(ReqT message) {
                    try (Tracer.SpanInScope ignored = tracer.withSpan(span)) {
                        super.onMessage(message);
                    }
                }

                @Override
                public void onHalfClose() {
                    try (Tracer.SpanInScope ignored = tracer.withSpan(span)) {
                        super.onHalfClose();
                    }
                }

                @Override
                public void onCancel() {
                    span.tag("rpc.cancelled", "true");
                    try (Tracer.SpanInScope ignored = tracer.withSpan(span)) {
                        super.onCancel();
                    } finally {
                        endOnce(span, ended);
                    }
                }
            };
        } catch (RuntimeException e) {
            span.error(e);
            endOnce(span, ended);
            throw e;
        }
    }

    private static void endOnce(Span span, AtomicBoolean ended) {
        if (ended.compareAndSet(false, true)) {
            span.end();
        }
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
