// jQuery plugin to perform easy and flexible form validations
// @author Ben Cessa <ben@pixative.com>

;(function( $ ){
  // Enable strict mode
  'use strict';

  // Namespace
  if( ! $.pixative ) {
    $.pixative = {};
  }

  // Object structure
  $.pixative.validator = function( element, options ) {
    // Self reference holder
    var self = this;

    // jQuery and DOM element references
    self.$el = $( element );
    self.el  = element;

    // Base regular expressions catalog, can be extended by the user by
    // adding or replacing entries in: options.additionalExpressions.*
    // jshint maxlen:200
    // Some of these expressions can be quite long
    self.expressions = {
      email: /^([\w\.\-]+)@([\w\-]+)((\.(\w){2,3})+)$/,
      alpha: /^[a-zA-Z]+$/,
      integer: /^(-?[1-9]\d*|0)$/,
      number: /^-?(?:0$0(?=\d*\.)|[1-9]|0)\d*(\.\d+)?$/,
      hex: /^#?([a-fA-F0-9]{6}|[a-fA-F0-9]{3})$/,
      postal: /^([1-9]{2}|[0-9][1-9]|[1-9][0-9])[0-9]{3}$/,
      phone: /^(\+\d{1,2}\s)?\(?\d{3}\)?[\s.-]\d{3}[\s.-]\d{4}$/,
      url: /^(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$/,
      cmc: /^5[1-5][0-9]{14}$/,
      cvisa: /^4[0-9]{12}(?:[0-9]{3})?$/,
      camex: /^3[47][0-9]{13}$/,
      cdiners: /^3(?:0[0-5]|[68][0-9])[0-9]{11}$/,
      cdisc: /^6(?:011|5[0-9]{2})[0-9]{12}$/,
      cjcb: /^(?:2131|1800|35\d{3})\d{11}$/,
      uuid4: /[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}/i,
      rfc: /^[A-Z]{3,4}[ |\-]{0,1}[0-9]{6}[ |\-]{0,1}[0-9A-Z]{3}$/i
    };

    // Base error messages catalog, can be extended by the user by
    // adding or replacing entries in: options.errorMessages.*
    self.messages = {
      required: 'This field is required',
      minlength: 'The minimum length for this field is %s, please verify',
      maxlength: 'The maximum length for this field is %s, please verify',
      email: 'The provided value doesn\'t appear to be a valid email address',
      alpha: 'Only characters or allowed for this field',
      integer: 'The provided value must be an integer number',
      number: 'The provided value must be a number',
      hex: 'The provided value must be a valid HEX value',
      postal: 'The provided value doesn\'t appear to be a valid postal code',
      phone: 'The provided value doesn\'t appear to be a valid phone number',
      url: 'The provided value doesn\'t appear to be a valid URL address',
      cmc: 'The provided value doesn\'t appear to be a valid MasterCard card',
      cvisa: 'The provided value doesn\'t appear to be a valid Visa card',
      camex: 'The provided value doesn\'t appear to be a valid AMEX card',
      cdiners: 'The provided value doesn\'t appear to be a valid Diners card',
      cdisc: 'The provided value doesn\'t appear to be a valid Discover card',
      cjcb: 'The provided value doesn\'t appear to be a valid JCB card',
      uuid4: 'The provided value doesn\'t appear to be a valid UUID',
      rfc: 'The provided value doesn\'t appear to be a valid RFC'
    };

    // Initialization stuff
    self.init = function() {
      // Extend default settings with passed in options
      self.settings = $.extend( true, $.pixative.validator.defaults, options );

      // Extend the base expressions catalog with user provided entries
      $.extend( self.expressions, self.settings.additionalExpressions );

      // Merge the default error message catalog with user provided data
      $.extend( self.messages, self.settings.errorMessages );

      // Get items to validate and init the error message holder
      self.items = self.$el.find( self.settings.elements );
      self.items.attr( self.settings.msgHolder, '' );

      // Prevent submit event ?
      if( self.settings.preventSubmit ) {
        self.$el.on( 'submit', function( e ) {
          // If validation fail, don't submit the form
          // if( ! self.isFormValid() ) {
          //   e.preventDefault();
          // }
          e.preventDefault();
          self.isFormValid();
        });
      }

      // Validate items on blur ?
      if( self.settings.validateOnBlur ) {
        self.items.each( function( index, item ) {
          // jshint unused: false
          // We're not using index but is part of the function declaration
          $( item ).on( 'blur', function( e ) {
            e.preventDefault();
            self.validateItem( this );
          });
        });
      }
    };

    // Validate the hole form
    self.isFormValid = function() {
      // Validate each item
      self.items.each( function( index, item ) {
        // jshint unused: false
        // We're not using index but is part of the function declaration
        self.validateItem( item );
      });

      // Check for errors
      var errors = self.items.filter( '.' + self.settings.errorClass );
      if( errors.length > 0 ) {
        if( self.settings.autoFocus ) {
          errors.first().focus();
        }
        return false;
      }

      // So far so good!
      // Send success event and trigger callback
      self.$el.trigger( 'validationSuccess' );
      self.settings.onSuccess.call( self.el );
      return true;
    };

    // Validate a specific form element
    self.validateItem = function( item ) {
      var el = $( item );

      // Cleanup previous error markers
      el.removeClass( self.settings.errorClass );
      el.attr( self.settings.msgHolder, '' );

      // Base validators
      var baseValidators = ['required','minlength','maxlength'];
      baseValidators.forEach( function( entry ) {
        if( el.attr( self.settings.dataPrefix + entry ) !== undefined ) {
          self[ entry + 'Validator']( item );
        }
      });

      // Expressions : Validates regular expressions catalog
      for( var k in self.expressions ) {
        if( el.attr( self.settings.dataPrefix + k ) !== undefined ) {
          if( ! self.expressions[ k ].test( el.val() ) ) {
            self.setError( item, self.messages[ k ] );
            return;
          }
        }
      }
    };

    // Validate if an element satisfy the attribute of 'required'
    self.requiredValidator = function( item ) {
      if( $( item ).val().trim() === '' ) {
        self.setError( item, self.messages.required );
        return false;
      }
      return true;
    };

    // Validate if an element is of a given minimum length
    self.minlengthValidator = function( item ) {
      var value = $( item ).attr( self.settings.dataPrefix + 'minlength' );
      if( $( item ).val().trim().length < value ) {
        self.setError( item, self.messages.minlength.replace('%s', value) );
        return false;
      }
      return true;
    };

    // Validate if an element is of a given maximum length
    self.maxlengthValidator = function( item ) {
      var value = $( item ).attr( self.settings.dataPrefix + 'maxlength' );
      if( $( item ).val().trim().length > value ) {
        self.setError( item, self.messages.maxlength.replace('%s', value) );
        return false;
      }
      return true;
    };

    // Set a validation error on a given form element
    self.setError = function( item, msg ) {
      // Don't override existing error messages, if any
      if( $( item ).attr( self.settings.msgHolder ) !== '' ) {
        return;
      }

      // Attach class and error message
      $( item ).addClass( self.settings.errorClass )
               .attr( self.settings.msgHolder, msg );

      // Send error event and trigger callback
      self.$el.trigger( 'validationError', { item: item, message: msg } );
      self.settings.onError.call( self.el );
    };

    // Self initialize
    self.init();
  };

  // Default settings
  $.pixative.validator.defaults = {
    autoFocus: false,
    dataPrefix: 'data-validator-',
    msgHolder: 'data-validator-error',
    errorClass: 'validator-error',
    elements: 'input, textarea, select',
    validateOnBlur: true,
    preventSubmit: true,
    additionalExpressions: {},
    errorMessages: {},
    onError: function(){},
    onSuccess: function(){}
  };

  // Plugin wrapper
  $.fn.pixativeFormValidator = function( options ) {
    // jshint nonew:false
    // This is the recommended jQuery plugin structure
    return this.each( function() {
      ( new $.pixative.validator( this, options ) );
    });
  };
})( jQuery );
