(($) ->
  $.badger = (@item, options) ->
    @options = $.extend({remote_style: false}, options)
    if @options['remote_style'] == true
      $('head').append("<link rel='stylesheet' type='text/css' href='http://ezcurdia.me/badger/badger.min.css' />")

    easyNumber = (num) =>
      fraction = num
      for x in ['', 'k', 'M', 'G', 'T', 'P', 'E', 'Z', 'Y']
        if fraction < 1000
          return "#{parseInt(fraction)}#{x}"
        else
          fraction /= Math.pow(10, 3)

    setBadges = (languages) =>
      languagesSorted = Object.keys(languages).sort (a,b) -> languages[b]-languages[a]
      for lang in languagesSorted
        @item.append("<span class='badge-#{lang.toLowerCase()}' title='more than #{easyNumber(languages[lang])} lines of #{lang} code'></span>")

    $.get("https://merithub.herokuapp.com/github/#{@options.username}/languages/").done (data) => setBadges(data)
    true

  $.fn.badger = (options) ->
    $.badger($(this), options)

  return
) jQuery
