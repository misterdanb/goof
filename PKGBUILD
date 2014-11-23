# Maintainer: danb <danielbusch1992@googlemail.com>

pkgname=goof
pkgver=20130718
pkgrel=1
pkgdesc="Simple one file transfer webserver"
arch=('x86_64' 'i686')
url="http://www.google.com"
license=('Beerware')
#depends=('go' 'git')
makedepends=('go' 'git')
options=('!strip' '!emptydirs')
#source=('go-mtpfs::git://github.com/hanwen/go-mtpfs.git')
md5sums=('SKIP')
_goofgit=http://github.com/misterdanb/goof
_goncursesurl=code.google.com/p/goncurses
_goflagsurl=github.com/jessevdk/go-flags

pkgver() {
  cd ./
  # Get the date of the last commit, in YYYYMMDD format
  #git log -1 --pretty=%cd --date=short | sed 's/-//g'
}

build() {
  cd "$srcdir"
  GOPATH="$srcdir" go get -v ${_goncursesurl}
  GOPATH="$srcdir" go get -v ${_goflagsurl}
  cd "$srcdir"
  git clone ${_goofgit}
  mkdir -p bin
  cd bin
  GOPATH="$srcdir" go build "$srcdir/goof/goof.go"
}

check() {
  source /etc/profile.d/go.sh
}

package() {
  source /etc/profile.d/go.sh
  mkdir -p "$pkgdir/$GOPATH"
  cp -R --preserve=timestamps "$srcdir"/{src,pkg} "$pkgdir/$GOPATH"

  mkdir -p "$pkgdir/usr"
  cp -R --preserve=timestamps "$srcdir/bin" "$pkgdir/usr"

  # Package license (if available)
  for f in LICENSE COPYING; do
    if [ -e "$srcdir/src/goof/$f" ]; then
      install -Dm644 "$srcdir/src/goof/$f" "$pkgdir/usr/share/licenses/goof/$f"
    fi
  done
}
