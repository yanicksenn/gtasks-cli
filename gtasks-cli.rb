class GtasksCli < Formula
  desc "A CLI for managing Google Tasks"
  homepage "https://github.com/yanicksenn/gtasks-cli"
  url "https://github.com/yanicksenn/gtasks-cli/archive/v0.1.0.tar.gz"
  sha256 "e2a51e5f8dcf4c756b902080d871d81489979f61be4293fea896c36e599ffb4f"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "."
  end

  test do
    system "#{bin}/gtasks", "--help"
  end
end
