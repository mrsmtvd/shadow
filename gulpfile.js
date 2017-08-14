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
    util           = require('gulp-util'),
    runSequence    = require('run-sequence');

var COMPONENTS = __dirname + '/components',
    VENDORS_DASHBOARD = COMPONENTS + '/dashboard/assets/vendors/',
    VENDORS_PROFILING = COMPONENTS + '/profiling/assets/vendors/',
    DEV_ENV = 'development';

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

gulp.task('default', ['watch']);

gulp.task('build', ['clean'], function () {
    runSequence(['frontend', 'backend']);
});

gulp.task('clean', function() {
    del.sync('server');
    del.sync(VENDORS_DASHBOARD + '/*');
    del.sync(VENDORS_PROFILING + '/*');
    del.sync(COMPONENTS + '/**/assets/css/*.min.css');
    del.sync(COMPONENTS + '/**/assets/js/*.min.js');
});

/**
 * Develop
 */
var
    server = null,
    serverBinName = 'server';

// TODO: watch glide.yaml

gulp.task('server:build', function() {
    var build = new Date().getTime();
    var args = ['-ldflags', '-X "main.build=' + build + '"', '-o', serverBinName, '-v', './examples/base/'];

    if (process.env.NODE_ENV === DEV_ENV) {
        args.unshift('-race');
    }

    args.unshift('build');

    result = child.spawnSync('go', args);

    if (result.status !== 0) {
        util.log(util.colors.red('Error during "go install": ' + result.stderr));
    }

    return result;
});

gulp.task('server:spawn', function() {
    if (server) {
        server.kill();
    }

    server = child.spawn('./' + serverBinName);

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
    ], ['compress-components']);

    gulp.watch([
        COMPONENTS + '/**/*.html'
    ], function() {
        if (process.env.NODE_ENV !== DEV_ENV) {
            runSequence(
                'bindata',
                'server:build',
                'server:spawn'
            );
        } else {
            runSequence(
                'server:spawn'
            );
        }
    });
    
    gulp.watch([
        COMPONENTS + '/**/assets/css/*.css',
        COMPONENTS + '/**/assets/js/*.js'
    ], function() {
        if (process.env.NODE_ENV !== DEV_ENV) {
            runSequence(
                'bindata',
                'server:build',
                'server:spawn'
            );
        } else {
            runSequence(
                'server:build', // for change app version
                'server:spawn'
            );
        }
    });

    // proto
    gulp.watch(['*/**/*.proto'], ['protobuf']);

    // go source
    gulp.watch(['*.go', '*/**/*.go', '!**/bindata_assetfs.go', '!**/bindata.go'], function() {
        runSequence(
            'server:build',
            'server:spawn'
        );
    }, 'server');
});

gulp.task('watch', function() {
    runSequence(
        'bindata',
        'server:build',
        function() {
            reload.listen();
            runSequence(
                'server:watch',
                'server:spawn'
            )
        }
    );
});

/**
 * Frontend
 */
gulp.task('compress-components', function() {
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
});

gulp.task('frontend', ['compress-components'], function() {
    // vendor
    gulp.src([
        'bower_components/bootstrap/dist/**/*.min.css',
        'bower_components/bootstrap/dist/**/*.min.js',
        'bower_components/bootstrap/dist/**/glyphicons-*',
        '!**/bootstrap-theme*.css'
    ])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/bootstrap'));

    gulp.src(['bower_components/jquery/dist/jquery.min.js'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/jquery/js'));

    gulp.src([
        'bower_components/font-awesome/**/*.min.css',
        'bower_components/font-awesome/**/fontawesome-*',
        'bower_components/font-awesome/**/FontAwesome.otf'
    ])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/font-awesome'));

    gulp.src(['bower_components/nprogress/nprogress.js'])
        .pipe(uglify({
            mangle: false
        }))
        .pipe(rename({
            suffix: '.min'
        }))
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/nprogress/js'));
    gulp.src(['bower_components/nprogress/nprogress.css'])
        .pipe(cleanCss())
        .pipe(rename({
            suffix: '.min'
        }))
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/nprogress/css'));

    gulp.src(['bower_components/fastclick/lib/*.js'])
        .pipe(uglify({
            mangle: false
        }))
        .pipe(rename({
            suffix: '.min'
        }))
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/fastclick/js'));

    gulp.src(['bower_components/autosize/dist/*.min.js'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/autosize/js'));

    gulp.src(['bower_components/validator/*.js'])
        .pipe(uglify({
            mangle: false
        }))
        .pipe(rename({
            suffix: '.min'
        }))
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/validator/js'));

    // waitMe
    gulp.src(['bower_components/waitMe/*.min.js'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/waitMe/js'));

    gulp.src(['bower_components/waitMe/*.min.css'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/waitMe/css'));

    // progressbar
    gulp.src(['bower_components/bootstrap-progressbar/*.min.js'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/bootstrap-progressbar/js'));

    gulp.src(['bower_components/bootstrap-progressbar/css/bootstrap-progressbar-3.3.4.min.css'])
        .pipe(rename('bootstrap-progressbar.min.css'))
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/bootstrap-progressbar/css'));

    // ichecked
    gulp.src(['bower_components/iCheck/*.min.js'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/icheck/js'));

    gulp.src(['bower_components/iCheck/skins/flat/green.css'])
        .pipe(cleanCss())
        .pipe(rename({
            suffix: '.min'
        }))
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/icheck/css'));

    gulp.src(['bower_components/iCheck/skins/flat/green*.png'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/icheck/css'));

    // switchery
    gulp.src(['bower_components/switchery/dist/*.min.js'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/switchery/js'));

    gulp.src(['bower_components/switchery/dist/*.min.css'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/switchery/css'));

    // flipclock
    gulp.src(['bower_components/flipclock/compiled/*.min.js'])
        .pipe(gulp.dest(VENDORS_PROFILING + '/flipclock/js'));

    gulp.src(['bower_components/flipclock/compiled/*.css'])
        .pipe(cleanCss())
        .pipe(rename({
            suffix: '.min'
        }))
        .pipe(gulp.dest(VENDORS_PROFILING + '/flipclock/css'));
    
    // wysiwyg
    /*
    gulp.src(['bower_components/bootstrap-wysiwyg/js/dist/bootstrap-wysiwyg.min.js'])
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/bootstrap-wysiwyg/js'));

    gulp.src(['bower_components/jquery.hotkeys/jquery.hotkeys.js'])
        .pipe(uglify({
            mangle: false
        }))
        .pipe(rename({
            suffix: '.min'
        }))
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/jquery.hotkey/js'));

    gulp.src(['bower_components/google-code-prettify/src/prettify.js'])
        .pipe(uglify({
            mangle: false
        }))
        .pipe(rename({
            suffix: '.min'
        }))
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/google-code-prettify/js'));

    gulp.src(['bower_components/google-code-prettify/bin/prettify.min.css'])
        .pipe(cleanCss())
        .pipe(rename({
            suffix: '.min'
        }))
        .pipe(gulp.dest(VENDORS_DASHBOARD + '/google-code-prettify/css'));
    */
});

/**
 * Backend
 */
gulp.task('backend', [], function() {
    runSequence(
        'golang',
        ['bindata', 'protobuf'],
        'easyjson'
    );
});

gulp.task('golang', function() {
    return gulp.src(__dirname + '/**/*.go')
        .pipe(exec('goimports -w <%= file.path %>'))
        .pipe(exec.reporter(execOptions));
});

gulp.task('bindata', function() {
    return gulp.src([
        COMPONENTS + '/*/templates',
        COMPONENTS + '/*/assets',
        COMPONENTS + '/*/public'
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
                    contents: new Buffer(folders.join('/... ') + '/...')
                };
            }
        }))
        .pipe(exec('cd <%= file.path %> && go-bindata-assetfs <%= options.debug %> -ignore="(\.DS_Store|\.gitignore)" -pkg=<%= options.path.basename(file.path) %> <%= file.contents %>', {
            path: path,
            debug: process.env.NODE_ENV === DEV_ENV ? '-debug' : ''
        }))
        .pipe(exec.reporter(execOptions));
});

gulp.task('protobuf', function() {
    return gulp.src(COMPONENTS + '/**/*.proto')
        .pipe(exec('protoc --proto_path=<%= options.path.dirname(file.path) %>  --go_out=plugins=grpc:<%= options.path.dirname(file.path) %> <%= file.path %>', {
            path: path
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