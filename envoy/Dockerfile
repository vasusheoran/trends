FROM envoyproxy/envoy:v1.18.4
COPY ./envoy.yaml /etc/envoy/envoy.yaml
# COPY ./proto.pb /protos/

CMD /usr/local/bin/envoy -c /etc/envoy/envoy.yaml -l trace --log-path /tmp/envoy_info.log

EXPOSE 9901