# vim: set tabstop=2 shiftwidth=2 softtabstop=2 expandtab:

class CloudKeys
  constructor: () ->
    @entities = []
    @version = ""
    @password = '' #todo replace with user password
    $('#pw').focus().keyup (evt) =>
      if evt.keyCode is 13
        @password = $('#pw').val()
        $('#loader').removeClass('hide')
        @fetchData()
        $('#newEntityLink').click =>
          @showForm()
        $('#editEncPWLink').click =>
          @showEditEncPWForm()

        $('#passwordRequest').addClass('hide')

        $('#search').keyup =>
          `var that = this`
          @limitItems(@getItems($('#search').val()))
          return
        $('#search').focus()
        $(window).keyup (evt) =>
          if evt.altKey is true and evt.keyCode is 66
            if typeof window.copyToClipboard is "function"
              copyToClipboard($('#items li.active .username').val())
            else
              $('#items li.active .username').focus().select()
          if evt.altKey is true and evt.keyCode is 79 # workaround to copy password very fast
            if typeof window.copyToClipboard is "function"
              copyToClipboard($('#items li.active .password').data('toggle'))
            else
              $('#items li.active .passwordtoggle em').click()
              $('#items li.active .password').focus().select()
          if evt.altKey is true and evt.keyCode is 80
            if typeof window.copyToClipboard is "function"
              copyToClipboard($('#items li.active .password').data('toggle'))
            else
              $('#items li.active .password').focus().select()
          if evt.altKey is true and evt.keyCode is 85
            if typeof window.copyToClipboard is "function"
              copyToClipboard($('#items li.active .url').val())
            else
              $('#items li.active .url').focus().select()

  import: (xml) ->
    parsedXML = $.parseXML(xml)

    for group in $(parsedXML).find('group')
      tag = $(group).find('>title').text()
      for entry in $(group).find('entry')
        e = $(entry)
        entity = {}
        entity['title'] = e.find('title').text()
        entity['username'] = e.find('username').text()
        entity['password'] = e.find('password').text()
        entity['url'] = e.find('url').text()
        entity['comment'] = e.find('comment').text()
        entity['tags'] = tag
        @entities.push(entity)
    @updateData =>
      $('#import').val('')
      $('#importLink').click()

  updateData: (callback) ->
    encrypted = @encrypt(JSON.stringify(@entities))
    hash = CryptoJS.SHA1(encrypted).toString()

    $.post 'ajax', {'version': @version, 'checksum': hash, 'data': encrypted}, (result) =>
      if result.error is true
        alert "An error occured, please reload and try it again"
      else
        if typeof callback isnt "undefined"
          callback()
        @updateInformation(result)
    , "json"

  fetchData: () ->
    $.get 'ajax', (data) =>
      @updateInformation(data)
    , "json"

  updateInformation: (data) ->
    @version = data.version

    if data.data == ""
      @entities = []
    else
      try
        @entities = $.parseJSON(@decrypt(data.data))
      catch e
        window.location.reload()

    @entities.sort(@sortItems)

    @showItems(@getItems(''))
    @limitItems(@getItems($('#search').val()))

  encrypt: (value) ->
    return CryptoJS.AES.encrypt(value, @password).toString()

  decrypt: (value) ->
    return CryptoJS.AES.decrypt(value, @password).toString(CryptoJS.enc.Utf8)

  getClipboardCode: (value) ->
    cb = $('<div class="clipboard"></div>')
    cb.click (e) ->
      elem = $("<textarea>#{ value }</textarea>").css({
        'position': 'absolute',
        'left': '-9999px',
        'readonly': 'readonly',
        'top': (window.pageYOffset || document.documentElement.scrollTop) + 'px'
      })

      $("body").append(elem)
      elem.focus()
      elem.select()
      document.execCommand('copy')
      elem.remove()
      return
    return cb

  limitItems: (items) ->
    $('#resultdescription span').text(items.length)
    current = 0
    $('#items>li').each (k, v) =>
      item = $(v)
      item.removeClass('odd')
      if $.inArray(item.data('num'), items) is -1
        item.addClass('hide')
      else
        if item.hasClass('hide')
          item.removeClass('hide')

        if current % 2 is 0
          item.addClass('odd')
        current = current + 1
      return
    return

  showItems: (items) ->
    $('#items li').remove()
    itemContainer = $('#items')
    $('#resultdescription span').text(items.length)
    for item, i in items
      additionalClass = ""
      if i % 2 is 0
        additionalClass = "odd"
      item = @entities[item]
      c = $("<li data-num=\"#{ item.num }\" class=\"#{ additionalClass }\">#{ item.title } <span>#{ item.username }</span></li>")
      ul = $("<ul></ul>")
      password = ""
      for char, i of item.password
        password += "*"

      field = $("<li><label>Username:</label><input type=\"text\" class=\"username\" value=\"#{ item.username }\"><br></li>")
      ul.append(field)
      @getClipboardCode(item.username).insertBefore(field.find("br"))

      field = $("<li class=\"passwordtoggle\"><label>Password:</label><input type=\"text\" class=\"password\" value=\"#{ password }\" data-toggle=\"#{ item.password.replace(/"/g, '&quot;') }\"><em> (toggle visibility)</em></span><br></li>")
      ul.append(field)
      @getClipboardCode(item.password).insertBefore(field.find("br"))

      field = $("<li><label>URL:</label><input type=\"text\" class=\"url\" value=\"#{ item.url }\"><br></li>")
      ul.append(field)
      @getClipboardCode(item.url).insertBefore(field.find("br"))

      lines_match = item.comment.match(/\n/g)
      if lines_match isnt null
        counter = lines_match.length
      if counter < 2
        counter = 2

      field = $("<li><label>Comment:</label><textarea class=\"comment\" rows=\"#{ counter + 2 }\">#{ item.comment }</textarea><br></li>")
      ul.append(field)
      @getClipboardCode(item.comment).insertBefore(field.find("br"))

      field = $("<li><label>Tags:</label><input type=\"text\" class=\"tags\" value=\"#{ item.tags }\"><br></li>")
      ul.append(field)
      @getClipboardCode(item.tags).insertBefore(field.find("br"))

      ul.append("<li class=\"last\"><button class=\"btn btn-primary\">Edit</button><br></li>")
      ul.find('.btn-primary').click (e) =>
        `var t = e.currentTarget`
        num = $(t).parent().parent().parent().data('num')
        if typeof num isnt "undefined" and typeof num isnt null
          @showForm(num)

      ul.find('.passwordtoggle em').click (e) =>
        `var t = e.currentTarget`
        elem = $(t).parent().find('.password')
        original = elem.data('toggle')
        elem.data('toggle', elem.val())
        elem.val(original)
      c.append(ul)

      c.click (e) =>
        `var that = this`
        elem = $(e.currentTarget)
        if elem.hasClass('active') is false
          $('#items li.active').removeClass('active').find('ul').slideUp()
          elem.addClass('active')
          elem.find('ul').slideDown()

      c.find('input').focus().select()

      itemContainer.append(c)

    $('.hide').removeClass('hide')
    $('#loader').addClass('hide')
    $('#passwordRequest').addClass('hide')
    $('#search').focus()

    return

  getItems: (search) ->
    result = []
    search = search.toLowerCase()
    for item, i in @entities
      if item.title.toLowerCase().indexOf(search) != -1 or item.username.toLowerCase().indexOf(search) != -1 or item.tags.toLowerCase().indexOf(search) != -1
        item.num = i
        result.push(i)

    return result

  sortItems: (a, b) ->
    aTitle = a.title.toLowerCase()
    bTitle = b.title.toLowerCase()
    `((aTitle < bTitle) ? -1 : ((aTitle > bTitle) ? 1 : 0))`

  showForm: (num) ->
    $('#editDialog input').val('')
    $('#editDialog textarea').val('')
    $('#editDialog .hide').removeClass('hide')
    fields = ['title', 'username', 'password', 'url', 'comment', 'tags']

    if typeof num isnt "undefined" and typeof @entities[num] isnt "undefined"
      $('#editDialog input[name="num"]').val(num)
      for elem in fields
        $("#editDialog ##{elem}").val(@entities[num][elem])
      $("#editDialog input#repeat_password").val(@entities[num]['password'])
    else
      $('#editDialog button.btn-danger').addClass('hide')

    $('#editDialog').modal({})
    $('#editDialog .btn-danger').unbind('click').click =>
      confirmation = confirm('Are you sure?')
      if confirmation is true
        num = $('#editDialog input[name="num"]').val()
        if typeof num isnt "undefined" and typeof num isnt null and num != ""
          @entities.splice(num, 1)

          @updateData =>
            $('#formClose').click()

    $('#editDialog .btn-primary').unbind('click').click =>
      if @validateForm()
        num = $('#editDialog input[name="num"]').val()
        entity = {}
        for field in fields
          entity[field] = $("##{field}").val()
        if typeof num != "undefined" and num != ""
          @entities[num] = entity
        else
          @entities.push(entity)

        @updateData =>
          $('#formClose').click()
      return

  validateForm: () ->
    $('#editDialog .has-error').removeClass('has-error')
    success = true
    if $('#title').val() == ""
      $('#title').parent().addClass('has-error')
      success = false

    if $('#password').val() isnt "" && $('#repeat_password').val() isnt $('#password').val()
      $('#password, #repeat_password').parent().addClass('has-error')
      success = false

    return success

  showEditEncPWForm: () ->
    $('#editEncPWDialog input').val('')
    $('#editEncPWDialog .hide').removeClass('hide')
    $('#editEncPWDialog').modal({})

    $('#editEncPWDialog .btn-primary').unbind('click').click =>
      $('#editEncPWDialog .has-error').removeClass('has-error')
      confirmation = confirm('Do you really want to update your encryption password?')

      if confirmation isnt true
        return

      if $('#editEncPW_current_password').val() isnt @password
        $('#editEncPW_current_password').parent().addClass('has-error')
        return

      new_password = $('#editEncPW_password').val()

      if new_password is undefined or new_password is ''
        $('#editEncPW_password').parent().addClass('has-error')
        return

      if new_password isnt $('#editEncPW_repeat_password').val()
        $('#editEncPW_password, #editEncPW_repeat_password').parent().addClass('has-error')
        return

      @password = new_password
      @updateData =>
        $('#formEncPWClose').click()
        alert 'Your encryption password has been changed. Keep this in mind for later.'
      return


window.CloudKeys = new CloudKeys()
$('#importLink').click =>
  $('#importContainer').toggle(500)

$('#importContainer button').click =>
  window.CloudKeys.import($('#import').val())
  #$('#import').val('')
