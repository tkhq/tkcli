class Turnkey < Formula
  desc "Turnkey CLI"
  homepage "https://github.com/tkhq/tkcli"
  version "v1.0.4"
  license "Apache License 2.0"

  if Hardware::CPU.arm?
    url "https://github.com/tkhq/tkcli/raw/v1.0.4/dist/turnkey.darwin-aarch64", using: CurlDownloadStrategy
    sha256 "fdd30cffebefb58b113526ae0dd26e5b906910a83d8a5741eadb1e5e8f46aea0"

    def install
      bin.install "turnkey.darwin-aarch64" => "turnkey"
    end
  end
  if Hardware::CPU.intel?
    url "https://github.com/tkhq/tkcli/raw/v1.0.4/dist/turnkey.darwin-x86_64", using: CurlDownloadStrategy
    sha256 "51bf14f44f2f0082ebe3af846cd318dd59518c83426705f5792b820c7581de9d"

    def install
      bin.install "turnkey.darwin-x86_64" => "turnkey"
    end
  end

end
