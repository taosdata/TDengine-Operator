#!/bin/bash
wget https://github.com/BurntSushi/ripgrep/releases/download/13.0.0/ripgrep-13.0.0-x86_64-unknown-linux-musl.tar.gz
tar xvf ripgrep-13.0.0-x86_64-unknown-linux-musl.tar.gz
cp ripgrep-13.0.0*/rg ./rg -f
rm ripgrep-13.0.0-x86_64-unknown-linux-musl.tar.gz
rm ripgrep-13.0.0-x86_64-unknown-linux-musl -rf

wget https://github.com/lotabout/rargs/releases/download/v0.3.0/rargs-v0.3.0-x86_64-unknown-linux-gnu.tar.gz
tar xvf rargs-v0.3.0-x86_64-unknown-linux-gnu.tar.gz
rm -rf rargs-v0.3.0-x86_64-unknown-linux-gnu.tar.gz
