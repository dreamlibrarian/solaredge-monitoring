FROM public.ecr.aws/lambda/provided:al2 as build

ARG lambda_bin

COPY ${lambda_bin} /main

ENTRYPOINT ["/main"]