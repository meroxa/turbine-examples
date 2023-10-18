package com.meroxa;

import com.meroxa.turbine.Runner;
import com.meroxa.turbine.Turbine;
import com.meroxa.turbine.TurbineApp;
import com.meroxa.turbine.TurbineRecord;

import java.util.List;

import static java.time.LocalDateTime.now;

public class Main implements TurbineApp {
    public static void main(String[] args) {
        Runner.start(new Main());
    }

    @Override
    public void setup(Turbine turbine) {
        turbine
            .resource("source_name")
            .read("collection_name", null)
            .process(this::process)
            .writeTo(turbine.resource("a-mysql-destination"), "collection_enriched", null);
    }

    private List<TurbineRecord> process(List<TurbineRecord> records) {
        return records.stream()
            .map(r -> {
                var copy = r.copy();
                copy.setPayload("This was processed at " + now());

                return copy;
            })
            .toList();
    }
}
