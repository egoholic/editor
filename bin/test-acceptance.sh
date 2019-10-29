#!/bin/bash
# Runs acceptance tests
DBNAME=stoa_blogging_test_acceptance PORT=8000 DOMAIN=wearestoa.com godog --random
