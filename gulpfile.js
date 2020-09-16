var child          = require('child_process'),
    del            = require('del'),
    path           = require('path'),
    
    gulp           = require('gulp'),
    cleanCss       = require('gulp-clean-css'),
    env            = require('gulp-env'),
    exec           = require('gulp-exec'),
    filterBy       = require('gulp-filter-by'),
    groupAggregate = require('gulp-group-aggregate'),
    reload         = require('gulp-livereload'),
    rename         = require('gulp-rename'),
    uglify         = require('gulp-uglify'),
    util           = require('gulp-util');

var COMPONENTS = __dirname + '/components',
    VENDORS_DATABASE = COMPONENTS + '/database/internal/assets/vendors/',
    VENDORS_DASHBOARD = COMPONENTS + '/dashboard/internal/assets/vendors/',
    VENDORS_I18N = COMPONENTS + '/i18n/internal/assets/vendors/',
    VENDORS_PROFILING = COMPONENTS + '/profiling/internal/assets/vendors/',
    DEV_ENV = 'development',
    SERVER_BIN_NAME = 'server';

var execOptions = {
    err: true,
    stderr: true,
    stdout: true
};

// set env variables
//if (process.env.NODE_ENV === DEV_ENV) {
    env({
        file: __dirname + '/env.json'
    });
//}

gulp.task('clean', function(done) {
    del.sync(SERVER_BIN_NAME);
    del.sync(VENDORS_DASHBOARD + '/*');
    del.sync(VENDORS_PROFILING + '/*');
    del.sync(COMPONENTS + '/**/assets/css/*.min.css');
    del.sync(COMPONENTS + '/**/assets/js/*.min.js');

    done();
});

/**
 * Frontend
 */
gulp.task('compress-components', function(done) {
    gulp.src([COMPONENTS + '/**/*.css', '!**/*min.css'])
        .pipe(cleanCss())
        .pipe(rename({
            suffix: '.min'
        }))
        .pipe(gulp.dest(COMPONENTS));
    gulp.src([COMPONENTS + '/**/*.js', '!**/*min.js'])
        .pipe(uglify({
            mangle: false
        }))
        .pipe(rename({
            suffix: '.min'
        }))
        .pipe(gulp.dest(COMPONENTS));

    done();
});

gulp.task('frontend', gulp.series('compress-components', function(done) {
    /**
     * Vendors
     */

    // jquery
    gulp.src(['bower_components/jquery/dist/jquery.min.js'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/jquery/js'));

    // bootstrap
    gulp.src([
        'bower_components/bootstrap/dist/**/*.min.css',
        'bower_components/bootstrap/dist/**/*.min.js',
        'bower_components/bootstrap/dist/**/glyphicons-*',
        '!**/bootstrap-theme*.css'
    ])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/bootstrap'));

    // autosize
    gulp.src(['bower_components/autosize/dist/*.min.js'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/autosize/js'));

    // bootstrap-progressbar
    gulp.src(['bower_components/bootstrap-progressbar/*.min.js'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/bootstrap-progressbar/js'));

    gulp.src(['bower_components/bootstrap-progressbar/css/bootstrap-progressbar-3.3.4.min.css'])
        .pipe(rename('bootstrap-progressbar.min.css'))
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/bootstrap-progressbar/css'));

    // bootstrap-show-password
    gulp.src(['bower_components/bootstrap-show-password/dist/*password.min.js'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/bootstrap-show-password/js'));

    // datatables
    gulp.src(['bower_components/datatables.net/js/*.min.js'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/datatables.net/js'));

    gulp.src(['bower_components/datatables.net-bs/js/*.min.js'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/datatables.net-bs/js'));

    gulp.src(['bower_components/datatables.net-bs/css/*.min.css'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/datatables.net-bs/css'));

    gulp.src(['bower_components/datatables.net-fixedheader/js/*.min.js'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/datatables.net-fixedheader/js'));

    gulp.src(['bower_components/datatables.net-fixedheader-bs/css/*.min.css'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/datatables.net-fixedheader-bs/css'));

    gulp.src(['bower_components/datatables.net-responsive/js/*.min.js'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/datatables.net-responsive/js'));

    gulp.src(['bower_components/datatables.net-responsive-bs/js/*.min.js'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/datatables.net-responsive-bs/js'));

    gulp.src(['bower_components/datatables.net-responsive-bs/css/*.min.css'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/datatables.net-responsive-bs/css'));
    
    // echarts
    gulp.src(['bower_components/echarts/dist/echarts.common.min.js'])
        .pipe(rename('echarts.min.js'))
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/echarts/js'));

    // fastclick
    gulp.src(['bower_components/fastclick/lib/*.js'])
        .pipe(uglify({mangle:false}))
        .pipe(rename({suffix:'.min'}))
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/fastclick/js'));

    // flipclock
    gulp.src(['bower_components/flipclock/compiled/*.min.js'])
        .pipe(gulp.dest(VENDORS_PROFILING + '/flipclock/js'));

    gulp.src(['bower_components/flipclock/compiled/*.css'])
        .pipe(cleanCss())
        .pipe(rename({suffix:'.min'}))
        .pipe(gulp.dest(VENDORS_PROFILING + '/flipclock/css'));

    // font-awesome
    gulp.src(['bower_components/font-awesome/css/all.min.css'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/font-awesome/css'));
    gulp.src(['bower_components/font-awesome/webfonts/*'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/font-awesome/webfonts'));

    // highlightjs
    gulp.src(['bower_components/highlightjs/*.min.js'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/highlightjs/js'));

    gulp.src(['bower_components/highlightjs/styles/tomorrow.css'])
        .pipe(cleanCss())
        .pipe(rename({suffix:'.min'}))
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/highlightjs/css'));

    // iCheck
    gulp.src(['bower_components/iCheck/*.min.js'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/icheck/js'));

    gulp.src(['bower_components/iCheck/skins/flat/green.css'])
        .pipe(cleanCss())
        .pipe(rename({suffix:'.min'}))
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/icheck/css'));

    gulp.src(['bower_components/iCheck/skins/flat/green*.png'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/icheck/css'));

    // jquery.tagsinput
    gulp.src(['bower_components/jquery.tagsinput/src/*.js'])
        .pipe(uglify({mangle:false}))
        .pipe(rename({suffix:'.min'}))
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/jquery.tagsinput/js'));

    // nprogress
    gulp.src(['bower_components/nprogress/nprogress.js'])
        .pipe(uglify({mangle:false}))
        .pipe(rename({suffix:'.min'}))
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/nprogress/js'));
    gulp.src(['bower_components/nprogress/nprogress.css'])
        .pipe(cleanCss())
        .pipe(rename({suffix:'.min'}))
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/nprogress/css'));

    // pnotify
    gulp.src([
            'bower_components/pnotify/dist/pnotify.js',
            'bower_components/pnotify/dist/pnotify.buttons.js'
        ])
        .pipe(rename({suffix:'.min'}))
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/pnotify/js'));

    gulp.src([
            'bower_components/pnotify/dist/pnotify.css',
            'bower_components/pnotify/dist/pnotify.buttons.css'
        ])
        .pipe(rename({suffix:'.min'}))
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/pnotify/css'));

    // select2
    gulp.src(['bower_components/select2/dist/js/select2.min.js'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/select2/js'));
    
    gulp.src(['bower_components/select2/dist/css/*.min.css'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/select2/css'));

    // switchery
    gulp.src(['bower_components/switchery/dist/*.min.js'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/switchery/js'));

    gulp.src(['bower_components/switchery/dist/*.min.css'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/switchery/css'));

    // validator
    gulp.src(['bower_components/validator/*.js'])
        .pipe(uglify({mangle:false}))
        .pipe(rename({suffix:'.min'}))
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/validator/js'));

    // waitMe
    gulp.src(['bower_components/waitMe/*.min.js'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/waitMe/js'));

    gulp.src(['bower_components/waitMe/*.min.css'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/waitMe/css'));

    // bootstrap-languages
    gulp.src(['bower_components/bootstrap-languages/languages.css'])
        .pipe(cleanCss())
        .pipe(rename({suffix:'.min'}))
        .pipe(gulp.dest(VENDORS_I18N + '/bootstrap-languages'));
    gulp.src(['bower_components/bootstrap-languages/*.png'])
        .pipe(gulp.dest(VENDORS_I18N + '/bootstrap-languages'));

    // dropzonejs
    gulp.src(['bower_components/dropzonejs/dist/min/dropzone.min.js'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/dropzonejs/js'));

    gulp.src(['bower_components/dropzonejs/dist/min/dropzone.min.css'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/dropzonejs/css'));

    // bootstrap-daterangepicker
    gulp.src(['bower_components/bootstrap-daterangepicker/daterangepicker.css'])
        .pipe(cleanCss())
        .pipe(rename({suffix:'.min'}))
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/bootstrap-daterangepicker/css'));
    gulp.src(['bower_components/bootstrap-daterangepicker/daterangepicker.js'])
        .pipe(uglify({mangle:false}))
        .pipe(rename({suffix:'.min'}))
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/bootstrap-daterangepicker/js'));

    gulp.src(['bower_components/moment/min/*.min.js'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/moment/js'));

    // jQuery Smart Wizard
    gulp.src(['bower_components/jQuery-Smart-Wizard/js/jquery.smartWizard.js'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/jQuery-Smart-Wizard/js'));

    done();
}));

/**
 * Backend
 */
gulp.task('golang', function() {
    return gulp.src([__dirname + '/**/**/*.go', '!' + __dirname + '/vendor/**', '!/**/bindata_assetfs.go'])
        .pipe(exec('goimports -w <%= file.path %>'))
        .pipe(exec.reporter(execOptions));
});

gulp.task('lint', function(cb) {
    child.exec('golangci-lint run -c .golangci.yml', function (err, stdout, stderr) {
        console.log(stdout);
        console.log(stderr);
        cb(err);
    });
});

gulp.task('i18n', function() {
    return gulp.src([__dirname + '/**/*.po', '!' + __dirname + '/vendor/**'])
        .pipe(exec('msgfmt <%= file.path %> -o <%= options.path.dirname(file.path) %>/<%= options.path.basename(file.path, ".po") %>.mo', {
            path: path
        }))
        .pipe(exec.reporter(execOptions));
});

gulp.task('bindata', function() {
    var ignores = [
        '[.]DS_Store',
        '[.]gitignore',
        '.*?[.]go$',
    ];

    if (process.env.NODE_ENV !== DEV_ENV) {
        ignores.push('.*?(([^n]|^)|([^i]|^)n|([^m]|^)in|([^.]|^|^[.])min)[.][jJ][sS]');
        ignores.push('.*?(([^n]|^)|([^i]|^)n|([^m]|^)in|([^.]|^|^[.])min)[.][cC][sS][sS]');
        ignores.push('.*?[.]po$');
    }
    
    return gulp.src([
        'examples/**/templates',
        'examples/**/assets',
        'examples/**/migrations',
        'examples/**/locales',

        COMPONENTS + '/**/templates',
        COMPONENTS + '/**/assets',
        COMPONENTS + '/**/migrations',
        COMPONENTS + '/**/locales'
    ])
        .pipe(groupAggregate({
            group: function (file){
                return path.dirname(file.path);
            },
            aggregate: function (group, files){
                folders = [];
                
                for (var i in files) {
                    folders.push(files[i].path.replace(group + '/', ''))
                }
                
                return {
                    path: group,
                    contents: Buffer.from(folders.join('/... ') + '/...')
                };
            }
        }))
        .pipe(exec('cd <%= file.path %> && go-bindata-assetfs <%= options.debug %> -ignore="('+ ignores.join('|') +')" -o ./bindata_assetfs.go -pkg=<%= options.path.basename(file.path) %> -nometadata -nomemcopy <%= file.contents %>', {
            path: path,
            debug: process.env.NODE_ENV === DEV_ENV ? '-debug' : ''
        }))
        .pipe(exec.reporter(execOptions));
});

gulp.task('protobuf', function() {
    return gulp.src(COMPONENTS + '/**/*.proto')
        .pipe(exec('protoc --proto_path=<%= options.path.dirname(file.path) %> --go_out=plugins=grpc,Mgoogle/protobuf/timestamp.proto=github.com/golang/protobuf/ptypes/timestamp,Mgoogle/protobuf/duration.proto=github.com/golang/protobuf/ptypes/duration:<%= options.path.normalize(options.root + "/../../..") %> <%= file.path %>', {
            path: path,
            root: __dirname
        }))
        .pipe(exec.reporter(execOptions));
});

gulp.task('easyjson', function() {
    return gulp.src(COMPONENTS + '/**/*.go')
        .pipe(filterBy(function (file) {
            return file.contents.toString().indexOf('easyjson:json') > -1
        }))
        .pipe(exec('easyjson <%= file.path %>'))
        .pipe(exec.reporter(execOptions));
});

gulp.task('enumer', function() {
    return gulp.src(__dirname + '/**/*.go')
        .pipe(filterBy(function (file) {
            return file.contents.toString().indexOf('enumer:json') > -1
        }))
        .pipe(exec('enumer -type=componentStatus -trimprefix=ComponentStatus -output=component_status_enumer.go -transform=snake <%= options.path.dirname(file.path) %>', {
            path: path
        }))
        .pipe(exec.reporter(execOptions));
});

gulp.task('backend', gulp.series(
    gulp.parallel('i18n', 'protobuf', 'easyjson', 'enumer'),
    'bindata', 'golang', 'lint'
));

/**
 * Develop
 */
var server = null;

// TODO: watch glide.yaml

gulp.task('server:build', function(done) {
    var build = new Date().getTime();
    var args = ['-ldflags', '-X "main.build=' + build + '"', '-o', SERVER_BIN_NAME, '-v', './examples/base/'];

    if (process.env.NODE_ENV === DEV_ENV) {
        args.unshift('-race');
    }

    args.unshift('build');

    result = child.spawnSync('go', args);
    if (result.status !== 0) {
        util.log(util.colors.red('Error during "go install": ' + result.stderr));
    }

    done();
});

gulp.task('server:spawn', function() {
    if (server) {
        server.kill();
    }

    server = child.spawn('./' + SERVER_BIN_NAME);

    server.stdout.once('data', function() {
        reload.reload('/');
    });

    server.stdout.on('data', function(data) {
        var lines = data.toString().split('\n');
        for (var l in lines) {
            if (lines[l].length) {
                util.log(lines[l]);
            }
        }
    });

    server.stderr.on('data', function(data) {
        process.stdout.write(data.toString());
    });
});

gulp.task('server:watch', function() {
    // templates
    gulp.watch([
        COMPONENTS + '/**/assets/css/*.css',
        COMPONENTS + '/**/assets/js/*.js',
        '!' + COMPONENTS + '/**/assets/css/*.min.css',
        '!' + COMPONENTS + '/**/assets/js/*.min.js'
    ],  gulp.series('compress-components'));

    gulp.watch([
        COMPONENTS + '/**/*.html'
    ], function() {
        if (process.env.NODE_ENV !== DEV_ENV) {
            gulp.parallel('bindata', 'server:build', 'server:spawn');
        } else {
            gulp.parallel('server:spawn');
        }
    });

    gulp.watch([
        COMPONENTS + '/**/assets/css/*.css',
        COMPONENTS + '/**/assets/js/*.js'
    ], function() {
        if (process.env.NODE_ENV !== DEV_ENV) {
            gulp.parallel('bindata', 'server:build', 'server:spawn');
        } else {
            gulp.parallel('server:build', 'server:spawn');
        }
    });

    // proto
    gulp.watch(['*/**/*.proto'],  gulp.series('protobuf'));

    // go source
    gulp.watch(['*.go', '*/**/*.go', '!**/bindata_assetfs.go', '!**/bindata.go'], gulp.parallel(
        'server:build', 'server:spawn'
    ), SERVER_BIN_NAME);
});

gulp.task('watch', gulp.parallel('bindata', 'server:build',
    function() {
        reload.listen();
        return gulp.parallel('server:watch', 'server:spawn')
    }
));

gulp.task('default', gulp.series('watch'));

gulp.task('build', gulp.series('clean', gulp.parallel(['frontend', 'backend'])));