@extends('master')
@section('content')
    <div class="row">
        <div class="col-lg-4 col-lg-offset-4">
            <form class="form-signin" method="POST" action="{{action('Auth\AuthController@postLogin')}}">
                {{ csrf_field() }}
                <h2 class="form-signin-heading text-center">Please sign in</h2>
                <label class="sr-only" for="email">Email address</label>
                <input type="email" autofocus="" required="" placeholder="Email address" class="form-control" id="email" name="email" value="{{ old('email') }}">
                <label class="sr-only" for="password">Password</label>
                <input type="password" required="" placeholder="Password" class="form-control" id="password" name="password">
                <div class="checkbox">
                    <label>
                        <input type="checkbox" value="true" name="remember"> Remember me
                    </label>
                </div>
                <button type="submit" class="btn btn-lg btn-primary btn-block">Sign in</button>
            </form>
            @if (count($errors) > 0)
                <div class="alert alert-danger alert-dismissible signin-alert" role="alert">
                    <button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                    <ul>
                        @foreach ($errors->all() as $error)
                            <li>{{ $error }}</li>
                        @endforeach
                    </ul>
                </div>
            @endif
        </div>
    </div>

@endsection