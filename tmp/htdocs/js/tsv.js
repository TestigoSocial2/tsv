var TSV = {
  // Basic setup
  init: function() {
    // Plugins
    $('[data-toggle="tooltip"]').tooltip()
    $('.carousel').carousel({
      interval: 8000,
      keyboard: false,
      pause: "hover",
    });
    moment.locale('es');

    // Set random hero image
    if( $('#hero') ) {
      var img = "url('images/hero_photo_"+ (Math.floor(Math.random() * 4) + 1) +".png')";
      $('#hero > div.photo').css( 'background-image', img );
    }

    // Pre-register form setup
    if( $('form#preregister').length ) {
      this.preRegisterSetup();
    }

    // Register form setup
    if( $('form#setup').length ) {
      this.registerSetup();
    }

    // Display main stats
    if( $('#highlights').length ) {
      this.statsSetup();
    }

    // Indicators filter setup
    if( $('form#filterForm').length ) {
      this.filterSetup();
    }
  },

  // Load main stats
  statsSetup: function() {
    var data = {}
    var charts = {}
    var ui = {
      totalContracts: $('span.totalContracts'),
      totalBudget: $('span#totalBudget'),
      totalAward: $('span#totalAward'),
      totalAmount: $('span#totalAmount'),
      firstDate: $('span#firstDate'),
      lastDate: $('span#lastDate'),
      description: $('span#orgDescription')
    }

    // $('.counter').counterUp({
    //   delay: 25,
    //   time: 2500
    // });

    // Dynamically set bucket used, default to 'gacm'
    var url = '/stats/gacm';
    if( getParameter('bucket') ) {
      url = '/stats/' + getParameter('bucket');
    }

    // Load stats
    $.ajax({
      type: "GET",
      url: url,
      success: function( res ) {
        data = JSON.parse(res);
      }
    }).done(function() {
      // Adjust labels
      ui.firstDate.text( moment( data.firstDate ).format('LL') );
      ui.lastDate.text( moment( data.lastDate ).format('LL') );
      ui.description.text( data.description );
      ui.totalContracts.text( data.contracts.total );
      ui.totalBudget.text((data.contracts.budget / 1000000).toLocaleString(undefined, {
        minimumFractionDigits: 1,
        maximumFractionDigits: 1
      }));
      ui.totalAward.text((data.contracts.awarded / 1000000).toLocaleString(undefined, {
        minimumFractionDigits: 1,
        maximumFractionDigits: 1
      }));
      ui.totalAmount.text( data.contracts.budget.toLocaleString(undefined, {
        minimumFractionDigits: 1,
        maximumFractionDigits: 1
      }));

      // Prepare chart data
      var directP = ((data.method['limited'].budget * 100) / data.contracts.budget).toFixed(2);
      var limitedP = ((data.method['selective'].budget * 100) / data.contracts.budget).toFixed(2);
      var publicP = ((data.method['open'].budget * 100) / data.contracts.budget).toFixed(2);
      var charts = {
        limited: {
          c: false,
          data: {
            labels: [
              'Adjudicación Directa (%)',
              'Total Contratado (%)'
            ],
            datasets: [
              {
                data: [directP, (100 - directP).toFixed(2)],
                backgroundColor: ["#CCB3FF", "#EEEEEE"],
              }
            ]
          }
        },
        selective: {
          c: false,
          data: {
            labels: [
              'Adjudicación Directa (%)',
              'Invitación a cuando menos 3 personas (%)',
              'Total Contratado (%)'
            ],
            datasets: [
              {
                data: [directP, limitedP, (100 - (Number(directP) + Number(limitedP))).toFixed(2)],
                backgroundColor: ["#CCB3FF", "#FF6384", "#EEEEEE"],
              }
            ]
          }
        },
        open: {
          c: false,
          data: {
            labels: [
              'Adjudicación Directa (%)',
              'Invitación a cuando menos 3 personas (%)',
              'Licitación Pública (%)'
            ],
            datasets: [
              {
                data: [directP, limitedP, publicP],
                backgroundColor: ["#CCB3FF", "#FF6384", "#7DE7CF"],
              }
            ]
          }
        }
      }

      // Configure data slider
      var slides = $('#content-slides');
      slides.on('slid.bs.carousel', function(){
        var active = slides.find('div.active');
        var method = active.data('section');
        if( method ) {
          active.find('span.contracts').text( data.method[method].total );
          active.find('span.amount').text( data.method[method].budget.toLocaleString({
            useGrouping: true
          }));

          if( !charts[method].c ) {
            charts[method].c = new Chart( active.find('canvas'), {
              type: "pie",
              data: charts[method].data,
              options: {
                responsive: true,
                responsiveAnimationDuration: 500,
                padding: 10
              }
            });
          }
        }
      });
    });
  },

  // Handler user registration process
  registerSetup: function() {
    var agencies = $('div.agency-grid div.item');
    var selectedAgencies = $('div.agency-grid input#selectedAgencies');
    agencies.click(function() {
      // Toggle state
      var a = $(this);
      var lbl = a.find('span.label');
      a.toggleClass('selected');
      if( lbl.text().toLowerCase() == 'seguir' ) {
        lbl.text( 'siguiendo' );
      } else {
        lbl.text( 'seguir' );
      }

      // Adjust result
      var res = [];
      agencies.each( function( i, v ) {
        if( $(v).hasClass( 'selected') ) {
          res.push( $(v).data('value') );
        }
      });
      selectedAgencies.val(JSON.stringify(res));
    });

    var projects = $('div.project-grid div.item');
    var selectedProjects = $('div.project-grid input#selectedProjects');
    projects.click(function() {
      // Toggle state
      var p = $(this);
      var lbl = p.find('span.label');
      p.toggleClass('selected');
      if( lbl.text().toLowerCase() == 'seguir' ) {
        lbl.text( 'siguiendo' );
      } else {
        lbl.text( 'seguir' );
      }

      // Adjust result
      var res = [];
      projects.each( function( i, v ) {
        if( $(v).hasClass( 'selected') ) {
          res.push( $(v).data('value') );
        }
      });
      selectedProjects.val(JSON.stringify(res));
    });

    var form = $('form#registerForm');
    form.pixativeFormValidator({
      msgHolder: 'title',
      errorMessages: {
        required: 'El campo es requerido',
        minlength: 'El valor proporcionado debe ser de al menos %s caracteres',
        maxlength: 'El valor proporcionado debe ser de máximo %s caracteres',
        email: 'El valor porporcionado no parece ser una dirección de correo valida',
        integer: 'El valor proporcionado debe ser un número',
        phone: 'El valor porporcionado no parece ser un número telefonico valido'
      },
      onError: function() {
        $('.validator-error').tooltip('destroy');
        $('.validator-error').tooltip();
      },
      onSuccess: function() {
        // Prepare data
        var data = {}
        form.serializeArray().forEach(function(el) {
          if( el.value == 'on' ) {
            el.value = true;
          }
          if( el.value == 'off' ) {
            el.value = false;
          }
          data[el.name] = el.value;
        });
        data.selectedAgencies = JSON.parse(data.selectedAgencies);
        data.selectedProjects = JSON.parse(data.selectedProjects);

        // Submit
        $.ajax({
          type: "POST",
          url: '/profile',
          data: {
            profile: JSON.stringify(data)
          },
          success: function( res ) {
            alert('Tu usuario ha quedado registrado exitosamente en Testigo Social Virtual 2.0');
          }
        })
      }
    });
  },

  // Pre-register process
  preRegisterSetup: function() {
    var regx = /^([\w\.\-]+)@([\w\-]+)((\.(\w){2,3})+)$/;
    var form = $('form#preregister');
    form.on( 'submit', function( e ) {
      e.preventDefault();
      var email = form.find('input').val();
      if( regx.test(email) ) {
        $.ajax({
          type: "POST",
          url: 'preregister',
          data: {
            email: email
          },
          success: function( res ) {
            alert("Gracias por tu interés en Testigo Social 2.0. Pronto te informaremos cómo podrás participar en las compras públicas mediante esta nueva plataforma.");
          }
        });
      } else {
        alert("La dirección proporcionada no es una dirección de correo electrónico valida, favor de verificar!");
      }
    });
  }
}

// Helper method to retrieve a GET variable
function getParameter(name, url) {
  if (!url) {
    url = window.location.href;
  }
  name = name.replace(/[\[\]]/g, "\\$&");
  var regex = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)"),
  results = regex.exec(url);
  if (!results) return false;
  if (!results[2]) return false;
  return decodeURIComponent(results[2].replace(/\+/g, " "));
}
