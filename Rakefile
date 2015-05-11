require 'sprockets/standalone'

Sprockets::Standalone::RakeTask.new(:assets) do |task, sprockets|
  task.assets   = %w(badger.min.js badger.min.css *.png *.svg *.woff *.eot *.zvgz *.ttf)
  task.sources  = %w(assets vendor/assets)
  task.output   = File.expand_path('../public', __FILE__)
  task.compress = false
  task.digest   = false

  sprockets.js_compressor  = :uglifier
  sprockets.css_compressor = :sass
end
