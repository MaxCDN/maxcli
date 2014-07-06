function run_tests {
# maxtail
##
  tool="./maxtail/maxtail.go"

  output=`go run $tool -h`
  assert_grep "echo '$output'" "Usage" "$tool: display usage"

  if test "$ALIAS" && test "$TOKEN" && test "$SECRET"
  then # Run functional tests.

    # run with a few filters
    output=`go run $tool -i 86400 -n -f nginx --status 200 --ssl nossl`

    refute_grep "echo '$output'" "Usage" "$tool: display usage"

    refute_grep "echo $output" " 404 " \
      "$tool: does not have status code 404"
  fi
}

