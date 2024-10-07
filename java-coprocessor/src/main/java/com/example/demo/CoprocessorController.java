package com.example.demo;

import com.example.demo.models.CoprocessorStage;
import com.example.demo.models.RouterPayload;
import com.example.demo.handlers.RouterRequestHandler;

import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class CoprocessorController {
    private final RouterRequestHandler routerRequestHandler;

    public CoprocessorController() {
        this.routerRequestHandler = new RouterRequestHandler();
    }

    @PostMapping("/")
    public RouterPayload StageSelect(@RequestBody RouterPayload request) {
        CoprocessorStage stage = request.getStage();

        if (stage == CoprocessorStage.ROUTER_REQUEST) {
            request = routerRequestHandler.handle(request);
        }

        return request;
    }
}
