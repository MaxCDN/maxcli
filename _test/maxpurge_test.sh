function run_tests {
# maxpurge
##
  tool="./maxpurge/maxpurge.go"

  output=`go run $tool -h`
  assert_grep "echo '$output'" "Usage" "$tool: display usage"

  output=`go run $tool -z ''`
  assert_grep "echo '$output'" "Usage" "$tool: display usage"
  assert_grep "echo '$output'" "Usage" \
    "$tool: requires valid zone"

  if test "$ALIAS" && test "$TOKEN" && test "$SECRET"
  then # Run functional tests.
    output=`go run $tool -z 123456 2> /dev/null`
    assert_grep "echo '$output'" "Purge failed after" \
      "$tool: fails on bad zone"
  fi
}
