@extends('master')
@section('content')
    @include('sidebar')
    <div class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
        @yield('main')
    </div>
@endsection