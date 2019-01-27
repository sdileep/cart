#!/usr/bin/env bash

# Create a new FD to allow us to print to the console live AND capture the output
exec 3>&1
output=$(glide install 2>&1 | tee /dev/fd/3; exit ${PIPESTATUS[0]})
status=$?
if [ ${status} -ne 0 ]; then
  exit ${status}
fi

# If ANY warnings happened, fail. Warnings can be things like:
# * Lock file may be out of date. Hash check of YAML failed. You may need to run 'update'
bad_output=$(echo "${output}" | grep "\[WARN\]")
if [ ! -z "${bad_output}" ]; then
  echo "Warnings occurred during 'glide install'. Please address the following:"
  echo "${bad_output}"
  exit 1
fi
