'use strict';

// configuration, adapt paths/folders to your project
var frontPath = './public/',
    bowerPath = './bower_components/',
    destPaths = {
        scripts: frontPath + 'js/',
        styles: frontPath + 'css/',
        images: frontPath + 'images/',
        fonts: frontPath + 'fonts/'
    };

// --------------------------------------------------------------

var gulp = require('gulp'),
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
    del = require('del');

// use "--production" option to minify everything
var inProduction = ('production' in gulputil.env),
    srcPaths = {
        scripts: [
            bowerPath + 'jquery/dist/jquery.min.js',
            bowerPath + 'semantic-ui/dist/semantic.min.js',
            'resources/assets/js/**/*.js'
        ],
        styles: [
            bowerPath + 'semantic-ui/dist/semantic.min.css',
            'resources/assets/less/*.less'
        ],
        fonts: [
            bowerPath + 'semantic-ui/dist/themes/default/assets/fonts/*.*',
            'resources/assets/fonts/*.*'
        ],
        images: [
            bowerPath + 'semantic-ui/dist/themes/default/assets/images/*.*',
            'resources/assets/images/**/*.*'
        ]
    };

// WARNING: removes files from folders (folders are kept)
gulp.task('prune', function (cb) {
    del([
        frontPath + 'js/*.min.js',
        frontPath + 'css/*.min.css',
        frontPath + 'fonts/*.*'
    ], { force: true });

    return cache.clearAll(cb);
});

// minify and copy all JS (except vendor scripts, sourcemaps are commented and basically useless)
gulp.task('scripts', function () {
    return gulp.src(srcPaths.scripts)
        .pipe(concat('app.min.js'))
        .pipe(gulpif(inProduction, uglify()))
        .pipe(gulp.dest(destPaths.scripts))
        ;
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
gulp.task('fonts', function() {
    return gulp.src(srcPaths.fonts)
        .pipe(newer(destPaths.fonts))
        .pipe(gulp.dest(destPaths.fonts))
        ;
});

// run phpunit tests
gulp.task('phpunit', function () {
    gulp.src('phpunit.xml')
        .pipe(gulpif(inProduction, phpunit()));
});

// watch for file changes
gulp.task('watch', function () {
    gulp.watch(srcPaths.scripts, ['scripts']);
    gulp.watch(srcPaths.styles, ['styles']);
    gulp.watch(srcPaths.images, ['images']);
});

// default task
gulp.task('default', ['scripts', 'styles', 'images', 'phpunit']);

/*
 * usage:
 *
 * gulp (to manually run all tasks)
 * gulp watch (for automatic awesomness while developing your js, css, less)
 * gulp --production (to minify all the things)
 *
 * separate processing:
 * gulp prune (removes all js, css fonts AND clear the internal cache)
 * gulp scripts (only JS)
 * gulp styles (only CSS and LESS)
 * gulp images (only images)
 * gulp fonts (only bootstrap fonts)
 * gulp phpunit (only awesome PHP tests)
 * 
 */