package com.example.demo.models;

import com.fasterxml.jackson.annotation.JsonEnumDefaultValue;
import com.fasterxml.jackson.annotation.JsonValue;

public enum CoprocessorStage {
    ROUTER_REQUEST("RouterRequest"),
    ROUTER_RESPONSE("RouterResponse"),
    SUPERGRAPH_REQUEST("SupergraphRequest"),
    SUPERGRAPH_RESPONSE("SupergraphResponse"),
    SUBGRAPH_REQUEST("SubgraphRequest"),
    SUBGRAPH_RESPONSE("SubgraphResponse"),
    @JsonEnumDefaultValue
    UNKNOWN("Unknown");

    private String name;

    private CoprocessorStage(String name) {
        this.name = name;
    }

    @JsonValue
    public String getName() {
        return this.name;
    }

    public void setName(String name) {
        this.name = name;
    }
}
