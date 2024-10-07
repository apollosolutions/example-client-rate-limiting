package com.example.demo.models;

import com.fasterxml.jackson.annotation.JsonInclude;

@JsonInclude(JsonInclude.Include.NON_ABSENT)
class CoprocessorControl {
    private Integer breakVal;

    CoprocessorControl() {
    }

    CoprocessorControl(Integer breakVal) {
        this.breakVal = breakVal;
    }

    public Integer getBreak() {
        return this.breakVal;
    }

    public void setBreak(Integer breakVal) {
        this.breakVal = breakVal;
    }
}
