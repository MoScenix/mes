package com.team10.mes.ai.node.coder;

/** Commit/finish stage corresponding to node/coder. */
public final class CoderNode {
  @FunctionalInterface
  public interface Committer {
    void commit(String output) throws Exception;
  }

  public String run(String output, Committer committer) throws Exception {
    committer.commit(output);
    return output;
  }
}
