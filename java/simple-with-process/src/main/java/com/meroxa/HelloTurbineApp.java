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
            .fromSource("postgres", Map.of("foo", "bar"))
            .process(this::process)
            .toDestination("destination_name", Map.of("blah", "blah"));
    }

    private List<TurbineRecord> process(List<TurbineRecord> records) {
        return records.stream()
            .filter(r -> r.jsonGet("$.payload.after.id").equals(9582724))
            .map(r -> {
                var copy = r.copy();
                copy.setPayload("customer emails is" + copy.jsonGet("$.payload.after.customer_email"));

                return copy;
            })
            .toList();
    }
}