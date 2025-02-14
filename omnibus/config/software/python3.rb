name "python3"

default_version "3.11.5"

if ohai["platform"] != "windows"

  dependency "libxcrypt"
  dependency "libffi"
  dependency "ncurses"
  dependency "zlib"
  dependency ENV["OMNIBUS_OPENSSL_SOFTWARE"] || "openssl"
  dependency "bzip2"
  dependency "libsqlite3"
  dependency "liblzma"
  dependency "libyaml"

  source :url => "https://python.org/ftp/python/#{version}/Python-#{version}.tgz",
         :sha256 => "a12a0a013a30b846c786c010f2c19dd36b7298d888f7c4bd1581d90ce18b5e58"

  relative_path "Python-#{version}"

  python_configure_options = [
    "--with-ensurepip=yes" # We upgrade pip later, in the pip3 software definition
  ]

  if mac_os_x?
    python_configure_options.push("--enable-ipv6",
                          "--with-universal-archs=intel",
                          "--enable-shared")
  elsif linux_target?
    python_configure_options.push("--enable-shared",
                          "--enable-ipv6")
  elsif aix?
    # something here...
  end

  python_configure_options.push("--with-dbmliborder=")

  build do
    # 2.0 is the license version here, not the python version
    license "Python-2.0"

    env = with_standard_compiler_flags(with_embedded_path)
    configure(*python_configure_options, :env => env)
    command "make -j #{workers}", :env => env
    command "make install", :env => env

    # There exists no configure flag to tell Python to not compile readline support :(
    major, minor, bugfix = version.split(".")

    delete "#{install_dir}/embedded/lib/python#{major}.#{minor}/test"
    block do
      FileUtils.rm_f(Dir.glob("#{install_dir}/embedded/lib/python#{major}.#{minor}/lib-dynload/readline.*"))
      FileUtils.rm_f(Dir.glob("#{install_dir}/embedded/lib/python#{major}.#{minor}/distutils/command/wininst-*.exe"))
    end
  end

else
  dependency "vc_redist_14"

  # note that starting with 3.7.3 on Windows, the zip should be created without the built-in pip
  source :url => "https://dd-agent-omnibus.s3.amazonaws.com/python-windows-#{version}-amd64.zip",
         :sha256 => "1BD44AC628CF39C61FF037C715B05D5E9A2980C2049031F73CF291F90E95A0F3".downcase

  vcrt140_root = "#{Omnibus::Config.source_dir()}/vc_redist_140/expanded"
  build do
    # 2.0 is the license version here, not the python version
    license "Python-2.0"

    command "XCOPY /YEHIR *.* \"#{windows_safe_path(python_3_embedded)}\""
    command "copy /y \"#{windows_safe_path(vcrt140_root)}\\*.dll\" \"#{windows_safe_path(python_3_embedded)}\""

    # Install pip
    python = "#{windows_safe_path(python_3_embedded)}\\python.exe"
    command "#{python} -m ensurepip"
  end
end
