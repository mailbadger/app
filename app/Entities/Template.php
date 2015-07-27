<?php

namespace newsletters\Entities;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\SoftDeletes;
use Prettus\Repository\Contracts\Transformable;
use Prettus\Repository\Traits\TransformableTrait;

class Template extends Model implements Transformable
{
    use TransformableTrait, SoftDeletes;

    protected $table = 'templates';

    protected $fillable = [
        'name',
        'content',
    ];

    protected $dates = ['deleted_at'];

    public function campaigns()
    {
        return $this->hasMany('newsletters\Entities\Campaign');
	}
}
