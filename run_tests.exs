# run_tests.exs

defmodule TestRunner do
  @package_dir "package"

  def run do
    IO.puts("🔍 Finding directories in #{@package_dir}...")

    dirs_to_test =
      @package_dir
      |> File.ls!()
      |> Enum.map(&Path.join(@package_dir, &1))
      |> Enum.filter(&File.dir?/1)

    all_passed =
      Enum.reduce(dirs_to_test, true, fn dir, acc ->
        acc and run_go_test(dir)
      end)

    if all_passed do
      IO.puts("\n✅ All tests passed successfully.")
    else
      IO.puts("\n❌ Some tests failed. Please review the output above.")
      System.halt(1)
    end
  end

  defp run_go_test(dir) do
    IO.puts("\n🧪 Testing in #{dir}...")

    {output, exit_code} = System.cmd("go", ["test", "."], cd: dir, stderr_to_stdout: true)

    IO.puts(output)

    if exit_code == 0 do
      true
    else
      IO.puts("❌ Test failed in #{dir} with exit code: #{exit_code}")
      false
    end
  end
end

TestRunner.run()
