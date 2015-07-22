<?php
/**
 * Created by PhpStorm.
 * User: filip
 * Date: 15.7.15
 * Time: 18:44
 */

namespace newsletters\Http\Controllers;

class DashboardController extends Controller
{
    public function __construct()
    {
        $this->middleware('auth');
    }

    public function getIndex()
    {
        return view('dashboard.campaigns.list');
    }
}