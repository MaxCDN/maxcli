function run_tests {
# maxcurl
##
  tool="./maxcurl/maxcurl.go"

  output=`go run $tool -h`
  assert_grep "echo '$output'" "Usage" "$tool: display usage"

  output=`go run $tool`
  assert_grep "echo '$output'" "Usage" "$tool: display usage"
  assert_grep "echo '$output'" "missing path value" \
    "$tool: requires path"

  if test "$ALIAS" && test "$TOKEN" && test "$SECRET"
  then # Run functional tests.
    output=`go run $tool /account.json`
    assert_grep "echo '$output'" '"account":' \
      "$tool: has account node"
    refute_grep "echo '$output'" "Vary: Accept-Encoding" \
      "$tool: doesn't show headers w/o -i"

    output=`go run $tool -i /account.json`
    assert_grep "echo '$output'" "Vary: Accept-Encoding" \
      "$tool: -i shows headers"

  fi
}
