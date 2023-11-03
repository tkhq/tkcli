class Turnkey < Formula
  desc "Turnkey CLI"
  homepage "https://github.com/tkhq/tkcli"
  version ""
  license "Apache License 2.0"

  if Hardware::CPU.arm?
    url "https://github.com/tkhq/tkcli/raw//dist/turnkey.darwin-aarch64", using: CurlDownloadStrategy
    sha256 "6db0e8288543cca3bafcca8d222c5a5fb0fa6a838baae3e541d58d0264bd9ec9"

    def install
      bin.install "turnkey.darwin-aarch64" => "turnkey"
    end
  end
  if Hardware::CPU.intel?
    url "https://github.com/tkhq/tkcli/raw//dist/turnkey.darwin-x86_64", using: CurlDownloadStrategy
    sha256 "68ffb9e5dbb0945b5fd53a2e2ccadfdba1f620720865c23e275fc159f70d30de"

    def install
      bin.install "turnkey.darwin-x86_64" => "turnkey"
    end
  end

end
