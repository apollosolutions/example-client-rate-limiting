package com.example.demo.models;

import com.fasterxml.jackson.annotation.JsonInclude;

import java.util.ArrayList;
import java.util.LinkedHashMap;

@JsonInclude(JsonInclude.Include.NON_ABSENT)
public class RouterPayload {
    private Object control;
    private LinkedHashMap<String, ArrayList<String>> headers;
    private String id;
    private String method;
    private String sdl;
    private CoprocessorStage stage;
    private Integer version;

    RouterPayload() {
        this.control = "continue";
        this.id = "";
    }

    public Object getControl() {
        return this.control;
    }

    public void setControl(Object control) {
        this.control = control;
    }

    public LinkedHashMap<String, ArrayList<String>> getHeaders() {
        return this.headers;
    }

    public void setHeaders(LinkedHashMap<String, ArrayList<String>> headers) {
        this.headers = headers;
    }

    public String getId() {
        return this.id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getMethod() {
        return this.method;
    }

    public void setMethod(String method) {
        this.method = method;
    }

    public String getSDL() {
        return this.sdl;
    }

    public void setSDL(String sdl) {
        this.sdl = sdl;
    }

    public CoprocessorStage getStage() {
        return this.stage;
    }

    public void setStage(CoprocessorStage stage) {
        this.stage = stage;
    }

    public Integer getVersion() {
        return this.version;
    }

    public void setVersion(Integer version) {
        this.version = version;
    }
}
