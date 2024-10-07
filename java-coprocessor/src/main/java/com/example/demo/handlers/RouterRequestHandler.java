package com.example.demo.handlers;

import com.example.demo.models.RouterPayload;

public class RouterRequestHandler {
    public RouterPayload handle(RouterPayload request) {
        // This is the object sent by the Router that you can act upon to update headers, context, auth claims, etc
        // If you update the "control" property from "Continue" to something like { "break": 400 }, it will terminate the request and return the specified HTTP error
        // See: https://www.apollographql.com/docs/router/customizations/coprocessor/
        // The object sent by the Router has been mapped to a `RouterRequestBody` object by the `@RequestBody` annotation in the `CoprocessorController` class
        System.out.println(request.getHeaders());
        System.out.println(request.getMethod());
        System.out.println(request.getSDL());
        System.out.println(request.getVersion());
        System.out.println(request.getStage());
        System.out.println(request.getControl());
        System.out.println(request.getId());

        return request;
    }
}
