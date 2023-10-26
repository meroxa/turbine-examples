package com.meroxa;

import java.util.List;
import java.util.Map;

import com.meroxa.turbine.Turbine;
import com.meroxa.turbine.TurbineApp;
import com.meroxa.turbine.TurbineRecord;
import jakarta.enterprise.context.ApplicationScoped;

@ApplicationScoped
public class HelloTurbineApp implements TurbineApp {

    @Override
    public void setup(Turbine turbine) {
        turbine
            .fromSource("kafka", Map.of("foo", "bar"))
            .process(this::process)
            .toDestination("kafka", Map.of("blah", "blah"));
    }

    private List<TurbineRecord> process(List<TurbineRecord> records) {
        return records.stream()
            .map(r -> {
                var copy = r.copy();
                String email = (String) copy.jsonGet("$.after.customer_email");
                copy.jsonSet("$.after.customer_email", email.toLowerCase() + ".subdomain");

                return copy;
            })
            .toList();
    }
}