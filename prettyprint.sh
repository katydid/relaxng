set -xe
find ./RelaxTestSuite/ -name '*s.rng'|while read fname; do
	dir=$(dirname "${fname}")
	cat $fname | xmllint --format - > $dir/p.rng
done