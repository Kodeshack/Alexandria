 #!/usr/bin/env sh

if [ ! -f $HOME/dep/sassc/bin/sassc ]
then
	mkdir -p $HOME/dep
	cd $HOME/dep
	rm -r libsass sassc
	git clone --branch 3.5.5 --depth 1 https://github.com/sass/libsass.git
	git clone --branch 3.5.0 --depth 1 https://github.com/sass/sassc.git
	SASS_LIBSASS_PATH=$(pwd)/libsass make -C sassc -j $(nproc)
fi
