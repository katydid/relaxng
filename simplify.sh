set -xe

# This script is used to simplify all the correct schemas in the Relax Test Suite.
# It outputs these simplified grammars to the respective folders where the correct schema was found as a s.rng file.

wget http://www.kohsuke.org/relaxng/rng2srng/rng2srng-20020831.zip
unzip rng2srng-20020831.zip
mv ./rng2srng-20020831/* .
find ./RelaxTestSuite/ -name '*c.rng'|while read fname; do
	echo $fname
	dir=$(dirname "${fname}")
	java -jar rng2srng.jar $fname > $dir/s.rng
done