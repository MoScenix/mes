package com.team10.mes.user.service;

public final class UnauthorizedException extends RuntimeException {
  public UnauthorizedException() {
    super("unauthorized");
  }
}
