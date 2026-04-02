project('altstash',
  version: '0.1.0',
  meson_version: '>= 0.64.0',
)

dependency('glib-2.0')
dependency('gtk4', version: '>= 4.10.0')
dependency('libadwaita-1', version: '>= 1.4')

go_bin = find_program('go')

custom_target('altstash',
  build_by_default: true,
  build_always_stale: true,
  output: 'altstash',
  console: true,
  install: true,
  install_dir: get_option('bindir'),
  command: [go_bin, 'build', '-v', '-o', '@OUTPUT@',
            meson.current_source_dir()]
)

subdir('data')
