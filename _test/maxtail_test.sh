function run_tests {
# maxtail
##
  tool="./maxtail/maxtail.go"

  output=`go run $tool -h`
  assert_grep "echo '$output'" "Usage" "$tool: display usage"

  if test "$ALIAS" && test "$TOKEN" && test "$SECRET"
  then # Run functional tests.
    refute_grep "go run $tool --no-follow -i 5" "Usage" \
      "$tool: doesn't display usage"
  fi
}

