'use strict';

/*
 * usage:
 *
 * gulp (to manually run all tasks)
 * gulp watch (for automatic awesomness while developing your js, css, less)
 * gulp --production (to minify all the things)
 *
 * separate processing:
 * gulp prune (removes all js, css fonts AND clear the internal cache)
 * gulp scripts
 * gulp styles
 * gulp images
 * gulp fonts
 *
 */

// configuration, adapt paths/folders to your project
var frontPath = './public/',
    npmPath = './node_modules/',
    destPaths = {
        scripts: frontPath + 'js/',
        styles: frontPath + 'css/',
        images: frontPath + 'images/',
        fonts: frontPath + 'fonts/'
    };

// --------------------------------------------------------------

var gulp = require('gulp'),
    browserify = require('browserify'),
    source = require('vinyl-source-stream'),
    buffer = require('vinyl-buffer'),
    rename = require('gulp-rename'),
    less = require('gulp-less'),
    concat = require('gulp-concat'),
    cache = require('gulp-cache'),
    notify = require('gulp-notify'),
    uglify = require('gulp-uglify'),
    minifycss = require('gulp-minify-css'),
    newer = require('gulp-newer'),
    imagemin = require('gulp-imagemin'),
    phpunit = require('gulp-phpunit'),
    gulputil = require('gulp-util'),
    gulpif = require('gulp-if'),
    del = require('del'),
    es = require('event-stream'),
    reactify = require('reactify');

// use "--production" option to minify everything
var inProduction = ('production' in gulputil.env),
    srcPaths = {
        scripts: [
            'resources/assets/js/global.js',
            'resources/assets/js/campaigns.jsx',
            'resources/assets/js/templates.jsx',
            'resources/assets/js/create-new-campaign.jsx',
            'resources/assets/js/create-new-template.jsx',
            'resources/assets/js/create-new-list.jsx',
            'resources/assets/js/subscribers.jsx',
            'resources/assets/js/reports.jsx',
            'resources/assets/js/settings.jsx'
        ],
        styles: [
            npmPath + 'select2/dist/css/select2.min.css',
            npmPath + 'magnific-popup/dist/magnific-popup.css',
            'resources/assets/less/*.less'
        ],
        fonts: [
            npmPath + 'bootstrap/fonts/*.*',
            'resources/assets/fonts/*.*'
        ],
        images: [
            'resources/assets/images/**/*.*'
        ]
    };

// WARNING: removes files from folders (folders are kept)
gulp.task('prune', function (cb) {
    del([
        frontPath + 'js/*.min.js',
        frontPath + 'css/*.min.css',
        frontPath + 'fonts/*.*'
    ], {force: true});

    return cache.clearAll(cb);
});

gulp.task('scripts', function () {
    // map them to our stream function
    var tasks = srcPaths.scripts.map(function (entry) {
        var dirnames = entry.substr(entry.indexOf('js') + 3).split('/');
        dirnames.pop();
        var dirname = dirnames.join('/');

        return browserify({
            entries: [entry],
            transform: [reactify]
        })
            .bundle()
            .pipe(source(entry))
            // rename them to have "bundle as postfix"
            .pipe(buffer())
            .pipe(gulpif(inProduction, uglify()))
            .pipe(rename({
                dirname: dirname,
                extname: '.bundle.js'
            }))
            .pipe(gulp.dest(destPaths.scripts));
    });
    // create a merged stream
    return es.merge.apply(null, tasks);
});

// css, less
gulp.task('styles', ['fonts'], function () {
    return gulp.src(srcPaths.styles)
        .pipe(less())
        .pipe(concat('app.min.css'))
        .pipe(gulpif(inProduction, minifycss()))
        .pipe(gulp.dest(destPaths.styles))
        ;
});

// optimize and copy only changed images
gulp.task('images', function () {
    return gulp.src(srcPaths.images)
        .pipe(newer(destPaths.images))
        .pipe(gulp.dest(destPaths.images))
        ;
});

// copy bootstrap and other newer fonts
gulp.task('fonts', function () {
    return gulp.src(srcPaths.fonts)
        .pipe(newer(destPaths.fonts))
        .pipe(gulp.dest(destPaths.fonts))
        ;
});

// watch for file changes
gulp.task('watch', function () {
    gulp.watch(srcPaths.scripts, ['scripts']);
    gulp.watch(srcPaths.styles, ['styles']);
    gulp.watch(srcPaths.images, ['images']);
});

// default task
gulp.task('default', ['scripts', 'styles', 'images']);
