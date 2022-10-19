FROM turbine-rb-base:0.0.1
WORKDIR /app
COPY . .
RUN bundle install
USER nobody
ENTRYPOINT ["./function_runner.sh"]