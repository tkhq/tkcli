class Turnkey < Formula
  desc "Turnkey CLI"
  homepage "https://github.com/tkhq/tkcli"
  version "v1.1.1"
  license "Apache License 2.0"

  if Hardware::CPU.arm?
    url "https://github.com/tkhq/tkcli/raw/v1.1.1/dist/turnkey.darwin-aarch64", using: CurlDownloadStrategy
    sha256 "e57e8abdaf069d8fc7574472140cf593d8a5ab20fe7e8cab2966daff2d1a54d7"

    def install
      bin.install "turnkey.darwin-aarch64" => "turnkey"
    end
  end
  if Hardware::CPU.intel?
    url "https://github.com/tkhq/tkcli/raw/v1.1.1/dist/turnkey.darwin-x86_64", using: CurlDownloadStrategy
    sha256 "bc644a9f0123a425982ff2e790b46151a9f2c5a0523e82d8a7d5ce26c547596b"

    def install
      bin.install "turnkey.darwin-x86_64" => "turnkey"
    end
  end

end
