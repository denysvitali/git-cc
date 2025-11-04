# Homebrew formula for git-cc
class GitCc < Formula
  desc "A conventional commit tool for git"
  homepage "https://github.com/denysvitali/git-cc"
  url "https://github.com/denysvitali/git-cc/archive/v{{.Version}}.tar.gz"
  sha256 "{{.TarballHash}}"
  license "MIT"

  depends_on "go"

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w -X main.version=v{{.Version}}")
    bin.install "git-cc"
  end

  test do
    system "#{bin}/git-cc --version"
  end
end