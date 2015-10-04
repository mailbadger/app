@extends('dashboard.layout')
@section('scripts')
    @parent
    <script src="//cdn.ckeditor.com/4.5.1/full/ckeditor.js"></script>
    <script type="text/javascript" src="{{asset('js/templates-list.bundle.js')}}"></script>
@endsection
@section('main')
    <h1 class="page-header">Templates</h1>
    <div id="templates"></div>
@endsection