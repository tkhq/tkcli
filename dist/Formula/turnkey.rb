class Turnkey < Formula
  desc "Turnkey CLI"
  homepage "https://github.com/tkhq/tkcli"
  version "v1.1.2"
  license "Apache License 2.0"

  if Hardware::CPU.arm?
    url "https://github.com/tkhq/tkcli/raw/v1.1.2/dist/turnkey.darwin-aarch64", using: CurlDownloadStrategy
    sha256 "4e6881ec55c1172aa6fcc870710eb4b36bcfc906cd3df334b4041c7ce558db32"

    def install
      bin.install "turnkey.darwin-aarch64" => "turnkey"
    end
  end
  if Hardware::CPU.intel?
    url "https://github.com/tkhq/tkcli/raw/v1.1.2/dist/turnkey.darwin-x86_64", using: CurlDownloadStrategy
    sha256 "118f9162601dd9c4a1f2f3630cb21ba63a832719596958ea2d416cff7fc9926a"

    def install
      bin.install "turnkey.darwin-x86_64" => "turnkey"
    end
  end

end
