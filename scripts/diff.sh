function diff_includes() {
  for i in $(git diff HEAD~ --name-only); do
  if [[ "$i" == $1* ]]; then
    return 0
  fi
  done

  return 1
}
