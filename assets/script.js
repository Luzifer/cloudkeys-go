// Generated by CoffeeScript 2.4.1
(function() {
  // vim: set tabstop=2 shiftwidth=2 softtabstop=2 expandtab:
  var CloudKeys;

  CloudKeys = class CloudKeys {
    constructor() {
      this.entities = [];
      this.version = "";
      this.password = ''; //todo replace with user password
      $('#pw').focus().keyup((evt) => {
        if (evt.keyCode === 13) {
          this.password = $('#pw').val();
          $('#loader').removeClass('hide');
          this.fetchData();
          $('#newEntityLink').click(() => {
            return this.showForm();
          });
          $('#editEncPWLink').click(() => {
            return this.showEditEncPWForm();
          });
          $('#passwordRequest').addClass('hide');
          $('#search').keyup(() => {
            this.limitItems(this.getItems($('#search').val()));
          });
          $('#search').focus();
          return $(window).keyup((evt) => {
            if (evt.altKey === true && evt.keyCode === 66) {
              if (typeof window.copyToClipboard === "function") {
                copyToClipboard($('#items li.active .username').val());
              } else {
                $('#items li.active .username').focus().select();
              }
            }
            if (evt.altKey === true && evt.keyCode === 79) { // workaround to copy password very fast
              if (typeof window.copyToClipboard === "function") {
                copyToClipboard($('#items li.active .password').data('toggle'));
              } else {
                $('#items li.active .passwordtoggle em').click();
                $('#items li.active .password').focus().select();
              }
            }
            if (evt.altKey === true && evt.keyCode === 80) {
              if (typeof window.copyToClipboard === "function") {
                copyToClipboard($('#items li.active .password').data('toggle'));
              } else {
                $('#items li.active .password').focus().select();
              }
            }
            if (evt.altKey === true && evt.keyCode === 85) {
              if (typeof window.copyToClipboard === "function") {
                return copyToClipboard($('#items li.active .url').val());
              } else {
                return $('#items li.active .url').focus().select();
              }
            }
          });
        }
      });
    }

    import(xml) {
      var e, entity, entry, group, j, l, len, len1, parsedXML, ref, ref1, tag;
      parsedXML = $.parseXML(xml);
      ref = $(parsedXML).find('group');
      for (j = 0, len = ref.length; j < len; j++) {
        group = ref[j];
        tag = $(group).find('>title').text();
        ref1 = $(group).find('entry');
        for (l = 0, len1 = ref1.length; l < len1; l++) {
          entry = ref1[l];
          e = $(entry);
          entity = {};
          entity['title'] = e.find('title').text();
          entity['username'] = e.find('username').text();
          entity['password'] = e.find('password').text();
          entity['url'] = e.find('url').text();
          entity['comment'] = e.find('comment').text();
          entity['tags'] = tag;
          this.entities.push(entity);
        }
      }
      return this.updateData(() => {
        $('#import').val('');
        return $('#importLink').click();
      });
    }

    updateData(callback) {
      var encrypted, hash;
      encrypted = this.encrypt(JSON.stringify(this.entities));
      hash = CryptoJS.SHA1(encrypted).toString();
      return $.post('ajax', {
        'version': this.version,
        'checksum': hash,
        'data': encrypted
      }, (result) => {
        if (result.error === true) {
          return alert("An error occured, please reload and try it again");
        } else {
          if (typeof callback !== "undefined") {
            callback();
          }
          return this.updateInformation(result);
        }
      }, "json");
    }

    fetchData() {
      return $.get('ajax', (data) => {
        return this.updateInformation(data);
      }, "json");
    }

    updateInformation(data) {
      var e;
      this.version = data.version;
      if (data.data === "") {
        this.entities = [];
      } else {
        try {
          this.entities = $.parseJSON(this.decrypt(data.data));
        } catch (error) {
          e = error;
          window.location.reload();
        }
      }
      this.entities.sort(this.sortItems);
      this.showItems(this.getItems(''));
      return this.limitItems(this.getItems($('#search').val()));
    }

    encrypt(value) {
      return CryptoJS.AES.encrypt(value, this.password).toString();
    }

    decrypt(value) {
      return CryptoJS.AES.decrypt(value, this.password).toString(CryptoJS.enc.Utf8);
    }

    getClipboardCode(value) {
      var cb;
      cb = $('<div class="clipboard" data-toggle="tooltip" data-original-title="Copied to clipboard!" data-trigger="manual"/>');
      cb.click(function(e) {
        var elem, t;
        elem = $(`<textarea>${value}</textarea>`).css({
          'position': 'absolute',
          'left': '-9999px',
          'readonly': 'readonly',
          'top': (window.pageYOffset || document.documentElement.scrollTop) + 'px'
        });
        $("body").append(elem);
        elem.focus();
        elem.select();
        document.execCommand('copy');
        elem.remove();
        t = $(this);
        t.tooltip('show');
        setTimeout(function() {
          return t.tooltip('hide');
        }, 1000);
      });
      return cb;
    }

    limitItems(items) {
      var current;
      $('#resultdescription span').text(items.length);
      current = 0;
      $('#items>li').each((k, v) => {
        var item;
        item = $(v);
        item.removeClass('odd');
        if ($.inArray(item.data('num'), items) === -1) {
          item.addClass('hide');
        } else {
          if (item.hasClass('hide')) {
            item.removeClass('hide');
          }
          if (current % 2 === 0) {
            item.addClass('odd');
          }
          current = current + 1;
        }
      });
    }

    showItems(items) {
      var additionalClass, c, char, counter, field, i, item, itemContainer, j, len, lines_match, password, ref, ul;
      $('#items li').remove();
      itemContainer = $('#items');
      $('#resultdescription span').text(items.length);
      for (i = j = 0, len = items.length; j < len; i = ++j) {
        item = items[i];
        additionalClass = "";
        if (i % 2 === 0) {
          additionalClass = "odd";
        }
        item = this.entities[item];
        c = $(`<li data-num="${item.num}" class="${additionalClass}">${item.title} <span>${item.username}</span></li>`);
        ul = $("<ul></ul>");
        password = "";
        ref = item.password;
        for (char in ref) {
          i = ref[char];
          password += "*";
        }
        field = $(`<li><label>Username:</label><input type="text" class="username" value="${item.username}"><br></li>`);
        ul.append(field);
        this.getClipboardCode(item.username).insertBefore(field.find("br"));
        field = $(`<li class="passwordtoggle"><label>Password:</label><input type="text" class="password" value="${password}" data-toggle="${item.password.replace(/"/g, '&quot;')}"><em> (toggle visibility)</em></span><br></li>`);
        ul.append(field);
        this.getClipboardCode(item.password).insertBefore(field.find("br"));
        field = $(`<li><label>URL:</label><input type="text" class="url" value="${item.url}"><br></li>`);
        ul.append(field);
        this.getClipboardCode(item.url).insertBefore(field.find("br"));
        lines_match = item.comment.match(/\n/g);
        if (lines_match !== null) {
          counter = lines_match.length;
        }
        if (counter < 2) {
          counter = 2;
        }
        field = $(`<li><label>Comment:</label><textarea class="comment" rows="${counter + 2}">${item.comment}</textarea><br></li>`);
        ul.append(field);
        this.getClipboardCode(item.comment).insertBefore(field.find("br"));
        field = $(`<li><label>Tags:</label><input type="text" class="tags" value="${item.tags}"><br></li>`);
        ul.append(field);
        this.getClipboardCode(item.tags).insertBefore(field.find("br"));
        ul.append("<li class=\"last\"><button class=\"btn btn-primary\">Edit</button><br></li>");
        ul.find('.btn-primary').click((e) => {
          var t = e.currentTarget;
          var num;
          num = $(t).parent().parent().parent().data('num');
          if (typeof num !== "undefined" && typeof num !== null) {
            return this.showForm(num);
          }
        });
        ul.find('.passwordtoggle em').click((e) => {
          var t = e.currentTarget;
          var elem, original;
          elem = $(t).parent().find('.password');
          original = elem.data('toggle');
          elem.data('toggle', elem.val());
          return elem.val(original);
        });
        c.append(ul);
        c.click((e) => {
          var elem;
          elem = $(e.currentTarget);
          if (elem.hasClass('active') === false) {
            $('#items li.active').removeClass('active').find('ul').slideUp();
            elem.addClass('active');
            return elem.find('ul').slideDown();
          }
        });
        c.find('input').focus().select();
        itemContainer.append(c);
      }
      $('.hide').removeClass('hide');
      $('#loader').addClass('hide');
      $('#passwordRequest').addClass('hide');
      $('#search').focus();
    }

    getItems(search) {
      var i, item, j, len, ref, result;
      result = [];
      search = search.toLowerCase();
      ref = this.entities;
      for (i = j = 0, len = ref.length; j < len; i = ++j) {
        item = ref[i];
        if (item.title.toLowerCase().indexOf(search) !== -1 || item.username.toLowerCase().indexOf(search) !== -1 || item.tags.toLowerCase().indexOf(search) !== -1) {
          item.num = i;
          result.push(i);
        }
      }
      return result;
    }

    sortItems(a, b) {
      var aTitle, bTitle;
      aTitle = a.title.toLowerCase();
      bTitle = b.title.toLowerCase();
      return ((aTitle < bTitle) ? -1 : ((aTitle > bTitle) ? 1 : 0));
    }

    showForm(num) {
      var elem, fields, j, len;
      $('#editDialog input').val('');
      $('#editDialog textarea').val('');
      $('#editDialog .hide').removeClass('hide');
      fields = ['title', 'username', 'password', 'url', 'comment', 'tags'];
      if (typeof num !== "undefined" && typeof this.entities[num] !== "undefined") {
        $('#editDialog input[name="num"]').val(num);
        for (j = 0, len = fields.length; j < len; j++) {
          elem = fields[j];
          $(`#editDialog #${elem}`).val(this.entities[num][elem]);
        }
        $("#editDialog input#repeat_password").val(this.entities[num]['password']);
      } else {
        $('#editDialog button.btn-danger').addClass('hide');
      }
      $('#editDialog').modal({});
      $('#editDialog .btn-danger').unbind('click').click(() => {
        var confirmation;
        confirmation = confirm('Are you sure?');
        if (confirmation === true) {
          num = $('#editDialog input[name="num"]').val();
          if (typeof num !== "undefined" && typeof num !== null && num !== "") {
            this.entities.splice(num, 1);
            return this.updateData(() => {
              return $('#formClose').click();
            });
          }
        }
      });
      return $('#editDialog .btn-primary').unbind('click').click(() => {
        var entity, field, l, len1;
        if (this.validateForm()) {
          num = $('#editDialog input[name="num"]').val();
          entity = {};
          for (l = 0, len1 = fields.length; l < len1; l++) {
            field = fields[l];
            entity[field] = $(`#${field}`).val();
          }
          if (typeof num !== "undefined" && num !== "") {
            this.entities[num] = entity;
          } else {
            this.entities.push(entity);
          }
          this.updateData(() => {
            return $('#formClose').click();
          });
        }
      });
    }

    validateForm() {
      var success;
      $('#editDialog .has-error').removeClass('has-error');
      success = true;
      if ($('#title').val() === "") {
        $('#title').parent().addClass('has-error');
        success = false;
      }
      if ($('#password').val() !== "" && $('#repeat_password').val() !== $('#password').val()) {
        $('#password, #repeat_password').parent().addClass('has-error');
        success = false;
      }
      return success;
    }

    showEditEncPWForm() {
      $('#editEncPWDialog input').val('');
      $('#editEncPWDialog .hide').removeClass('hide');
      $('#editEncPWDialog').modal({});
      return $('#editEncPWDialog .btn-primary').unbind('click').click(() => {
        var confirmation, new_password;
        $('#editEncPWDialog .has-error').removeClass('has-error');
        confirmation = confirm('Do you really want to update your encryption password?');
        if (confirmation !== true) {
          return;
        }
        if ($('#editEncPW_current_password').val() !== this.password) {
          $('#editEncPW_current_password').parent().addClass('has-error');
          return;
        }
        new_password = $('#editEncPW_password').val();
        if (new_password === void 0 || new_password === '') {
          $('#editEncPW_password').parent().addClass('has-error');
          return;
        }
        if (new_password !== $('#editEncPW_repeat_password').val()) {
          $('#editEncPW_password, #editEncPW_repeat_password').parent().addClass('has-error');
          return;
        }
        this.password = new_password;
        this.updateData(() => {
          $('#formEncPWClose').click();
          return alert('Your encryption password has been changed. Keep this in mind for later.');
        });
      });
    }

  };

  window.CloudKeys = new CloudKeys();

  $('#importLink').click(() => {
    return $('#importContainer').toggle(500);
  });

  $('#importContainer button').click(() => {
    return window.CloudKeys.import($('#import').val());
  });

  //$('#import').val('')

}).call(this);
